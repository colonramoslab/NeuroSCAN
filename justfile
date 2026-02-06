# List commands
default:
    @just --list

build-backend:
    @echo "Building backend..."
    go build -trimpath -buildmode=pie -buildvcs=false -ldflags "-w -s" -o ./neuroscan cmd/main.go

build-frontend:
    @echo "Building frontend..."
    cd frontend && NODE_OPTIONS=--openssl-legacy-provider npm run build

copy-vendored:
    @echo "Copying vendored assets..."
    cd frontend && cp ./overwrite/Canvas.js ./node_modules/@metacell/geppetto-meta-ui/3d-canvas/ && cp ./overwrite/ThreeDEngine.js ./node_modules/@metacell/geppetto-meta-ui/3d-canvas/threeDEngine/ && cp ./overwrite/MeshFactory.js ./node_modules/@metacell/geppetto-meta-ui/3d-canvas/threeDEngine/ && cp ./overwrite/OBJLoader.js ./node_modules/@metacell/geppetto-meta-ui/3d-canvas/threeDEngine/
    @echo "Done copying vendored assets."

build: build-backend build-frontend copy-vendored

test:
	@echo "Running tests..."
	go test -v ./...

modernize:
	@echo "Running gopls modernize..."
	go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -fix -test ./...

lint:
	@echo "Running golangci-lint..."
	golangci-lint run ./...

sec:
	@echo "Running gosec..."
	gosec ./...

vuln:
	@echo "Running gosec..."
	govulncheck ./...
