# ProxImage

Image proxy server that fetches remote images and converts them to WebP format.

## Quick Start

```bash
go run main.go
```

The server starts on `http://localhost:3000`.

## Usage

```
GET /image?url=<image_url>
```

**Example:**

```bash
curl "http://localhost:3000/image?url=https://example.com/photo.jpg" -o output.webp
```

**Response headers:**
- `Content-Type: image/webp`
- `X-Cache: HIT` or `X-Cache: MISS`
- `Cache-Control: public, max-age=86400`

## Features

- Converts JPEG, PNG, GIF, BMP, TIFF, WebP to WebP
- LRU in-memory cache (500 entries, 1 hour TTL)
- Rate limiting (30 requests/minute per IP)
- Max image size: 10 MB

## Requirements

- Go 1.25.0

Optionally, use the Nix dev shell:

```bash
nix develop
```

## Build

```bash
go build
```

## Test

```bash
go test
```

## Docker

```bash
docker compose up -d
```

The server runs on `http://localhost:3000`.

**Environment variables:**
- `GODEBUG=netdns=go+v4` - DNS resolution (Docker)