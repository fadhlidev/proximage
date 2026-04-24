package handler

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/fadhlidev/proximage/cache"

	"github.com/chai2010/webp"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

const (
	maxImageSize = 10 << 20 // 10 MB
	fetchTimeout = 10 * time.Second
	webpQuality  = 82.0
)

var httpClient = &http.Client{
	Timeout: fetchTimeout,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
	},
}

type ImageHandler struct {
	cache *cache.Cache
}

func New(c *cache.Cache) *ImageHandler {
	return &ImageHandler{cache: c}
}

func (h *ImageHandler) Convert(c fiber.Ctx) error {
	rawURL := c.Query("url")
	if rawURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "query param 'url' is required",
		})
	}

	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid url",
		})
	}

	if cached, ok := h.cache.Get(rawURL); ok {
		c.Set("Content-Type", "image/webp")
		c.Set("X-Cache", "HIT")
		return c.Send(cached)
	}

	webpData, err := fetchAndConvert(rawURL)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to process image: %v", err),
		})
	}

	h.cache.Set(rawURL, webpData)

	c.Set("Content-Type", "image/webp")
	c.Set("X-Cache", "MISS")
	c.Set("Cache-Control", "public, max-age=86400")
	return c.Send(webpData)
}

func fetchAndConvert(imageURL string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), fetchTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "WebP-Proxy/1.0")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upstream returned %d", resp.StatusCode)
	}

	limited := io.LimitReader(resp.Body, maxImageSize)
	raw, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}

	img, err := decodeImage(raw, resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}

	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, &webp.Options{Quality: webpQuality}); err != nil {
		return nil, fmt.Errorf("webp encode failed: %w", err)
	}

	return buf.Bytes(), nil
}

func decodeImage(data []byte, contentType string) (image.Image, error) {
	r := bytes.NewReader(data)

	ct := strings.ToLower(strings.Split(contentType, ";")[0])
	switch ct {
	case "image/jpeg", "image/jpg":
		return jpeg.Decode(r)
	case "image/png":
		return png.Decode(r)
	case "image/gif":
		return gif.Decode(r)
	case "image/bmp":
		return bmp.Decode(r)
	case "image/tiff":
		return tiff.Decode(r)
	case "image/webp":
		return webp.Decode(r)
	default:
		img, _, err := image.Decode(bytes.NewReader(data))
		return img, err
	}
}
