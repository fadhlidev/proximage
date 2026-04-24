# AGENTS.md

## Project Overview
Go web service using Fiber v3. Entry point: `main.go`

## Developer Commands
```bash
go run main.go    # Run dev server on :3000
go build         # Build binary
go test          # Run tests (none currently)
```

## Requirements
- Go 1.25.0
- Optional: Nix with flakes (`nix develop`)

## Key Facts
- HTTP server listens on port 3000
- No embedded templates or assets
- No database or external services
- No Codegen or build artifacts

## Commit Messages
Format: `[action]: [message]` (e.g., `add: new endpoint`, `fix: memory leak in cache`)