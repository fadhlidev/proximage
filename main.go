package main

import (
	"log"
	"runtime"
	"time"

	"github.com/fadhlidev/proximage/cache"
	"github.com/fadhlidev/proximage/handler"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	// Cache: max 500 entries, TTL 1 hour
	imgCache := cache.New(500, time.Hour)
	imgHandler := handler.New(imgCache)

	app := fiber.New(fiber.Config{
		BodyLimit:    1 * 1024 * 1024, // 1MB
		Concurrency:  runtime.NumCPU() * 256,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	})

	app.Use(recover.New())
	app.Use(logger.New())

	app.Use(limiter.New(limiter.Config{
		Max:        30,
		Expiration: time.Minute,
		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP()
		},
	}))

	app.Get("/image", imgHandler.Convert)

	log.Fatal(app.Listen(":3000", fiber.ListenConfig{
		EnablePrefork: true,
	}))
}
