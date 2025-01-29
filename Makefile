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
	cp .env /etc/default/neuroscan.env

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_DIR)/$(BINARY_NAME)

test:
	@echo "Running tests..."
	go test -v ./...
