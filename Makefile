.PHONY: build clean

# Binary name and path
BINARY_DIR=~/code/neuroscan
BINARY_NAME=neuroscan
CMD_PATH=cmd/neuroscan

# Build flags for smaller binary size
LDFLAGS=-w -s
BUILD_FLAGS=-trimpath -buildvcs=false -ldflags "$(LDFLAGS)"

build:
	@echo "Building binary..."
	go build $(BUILD_FLAGS) -o $(BINARY_DIR)/$(BINARY_NAME) ./$(CMD_PATH)

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_DIR)/$(BINARY_NAME)

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
