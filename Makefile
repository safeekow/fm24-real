.PHONY: build clean test install run-check run-apply run-update help

# Variables
BINARY_NAME=fm24-real
VERSION=1.0.0
BUILD_DIR=build
DIST_DIR=dist

# Build targets
build: ## Build the binary
	go build -o $(BINARY_NAME) -ldflags="-s -w" .

build-all: ## Build for all platforms
	@mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 -ldflags="-s -w" .
	GOOS=darwin GOARCH=arm64 go build -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 -ldflags="-s -w" .
	GOOS=windows GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe -ldflags="-s -w" .
	GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 -ldflags="-s -w" .

clean: ## Clean build artifacts
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)

test: ## Run tests
	go test -v ./...

install: build ## Install binary to /usr/local/bin
	install -m 755 $(BINARY_NAME) /usr/local/bin/

# Development shortcuts
run-check: build ## Run check command
	./$(BINARY_NAME) --check

run-apply: build ## Run apply command
	./$(BINARY_NAME) --apply

run-update: build ## Run update command
	./$(BINARY_NAME) --update

# Dependencies
deps: ## Download dependencies
	go mod download
	go mod tidy

# Help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
