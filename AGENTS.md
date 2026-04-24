# AGENTS.md

## Project Overview
Image proxy server that converts remote images to WebP format using Fiber v3.

## API
- `GET /image?url=<image_url>` - fetch and convert image to WebP

## Developer Commands
```bash
go run main.go    # Run dev server on :3000
go build          # Build binary
go test           # Run tests (none currently)
```

## Requirements
- Go 1.25.0
- Optional: Nix with flakes (`nix develop`)

## Architecture
- `main.go` - server setup, middleware, routing
- `handler/image.go` - image fetching, format detection, WebP conversion
- `cache/cache.go` - LRU in-memory cache

## Key Facts
- Cache: 500 entries, 1 hour TTL, X-Cache header (HIT/MISS)
- Rate limit: 30 req/min per IP
- Max image size: 10 MB, body limit: 1 MB
- Supports: JPEG, PNG, GIF, BMP, TIFF, WebP input → WebP output (quality 82)
- Prefork enabled for production

## Commit Messages
Format: `[action]: [message]` (e.g., `add: new endpoint`, `fix: memory leak in cache`)