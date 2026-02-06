# NeuroSCAN Development Guide

## Build/Test Commands
- **Build**: `make build` (creates optimized binary) or `go build -o neuroscan ./cmd/main.go`
- **Test**: `make test` or `go test -v ./...`
- **Single test**: `go test -v ./pkg/logging -run TestSpecificFunction`
- **Lint**: `make lint` (requires golangci-lint)
- **Security**: `make sec` (gosec) and `make vuln` (govulncheck)

## Architecture
- **Main**: CLI with 3 commands: `web` (API server), `ingest` (data ingestion), `transcode` (video processing)
- **Backend**: Go with Echo framework, PostgreSQL, Goose migrations, structured as domain/repository/service/handler
- **Frontend**: React with Redux Toolkit, Node 14.21.3, Yarn package manager, requires vendor file overwrites
- **Database**: PostgreSQL with ULIDs, migrations in `/migrations`, connection via `DB_DSN` env var
- **Data**: Ingests `.gltf` files from structured directories: `<STAGE>/<TIMEPOINT>/<CELL_TYPE>/<FILE>.gltf`

## Code Style (Go)
- **Imports**: stdlib → external → local (`neuroscan/`) with blank line separation
- **Naming**: PascalCase exports, camelCase private, constants with descriptive prefixes
- **Errors**: immediate `if err != nil` checks, wrap with `fmt.Errorf("%w", err)`
- **Structs**: JSON tags as snake_case, pointers for nullable fields, `ToDomain()` conversion methods
- **Dependencies**: inject via struct embedding, use interfaces for testability

## Environment
- Copy `.env.example` to `.env` and configure `DB_DSN`, `APP_FRONTEND_DIR`, `APP_GLTF_DIR`
- Frontend build required before serving: `yarn build` (env vars loaded from `frontend/.env.development`; override with `frontend/.env.development.local`)
