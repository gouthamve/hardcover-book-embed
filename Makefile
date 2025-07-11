.PHONY: all build run test clean dev help

# Variables
BINARY_NAME=hardcover-embed
GO=go
GOFLAGS=
SERVER_PATH=cmd/server/main.go

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) $(SERVER_PATH)
	@echo "Build complete: ./$(BINARY_NAME)"

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

# Run without building (assumes already built)
run-only:
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

# Development mode with auto-reload (requires air)
dev:
	@command -v air >/dev/null 2>&1 || { echo "Installing air..."; go install github.com/air-verse/air@latest; }
	@echo "Running in development mode with auto-reload..."
	air

# Run tests
test:
	@echo "Running Go tests..."
	$(GO) test -v ./...

# Run API integration tests (requires running server)
test-api:
	@echo "Running API integration tests..."
	@if [ -z "$(USERNAME)" ]; then \
		echo "Error: USERNAME is required"; \
		echo "Usage: make test-api USERNAME=your-username"; \
		exit 1; \
	fi
	@cd test && ./test.sh $(USERNAME)

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	$(GO) mod tidy

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -f test/test_hardcover
	@echo "Clean complete"

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Run linter (requires golangci-lint)
lint:
	@command -v golangci-lint >/dev/null 2>&1 || { echo "Installing golangci-lint..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; }
	@echo "Running linter..."
	golangci-lint run

# Check for vulnerabilities
vuln:
	@echo "Checking for vulnerabilities..."
	$(GO) list -json -m all | nancy sleuth

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GO) mod download

# Setup environment
setup: deps
	@echo "Setting up environment..."
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env file from .env.example"; \
		echo "Please edit .env and add your HARDCOVER_API_TOKEN"; \
	else \
		echo ".env file already exists"; \
	fi

# Start the server (production mode)
serve: build
	@echo "Starting server in production mode..."
	./$(BINARY_NAME)

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 $(GO) build -o dist/$(BINARY_NAME)-linux-amd64 $(SERVER_PATH)
	GOOS=darwin GOARCH=amd64 $(GO) build -o dist/$(BINARY_NAME)-darwin-amd64 $(SERVER_PATH)
	GOOS=darwin GOARCH=arm64 $(GO) build -o dist/$(BINARY_NAME)-darwin-arm64 $(SERVER_PATH)
	GOOS=windows GOARCH=amd64 $(GO) build -o dist/$(BINARY_NAME)-windows-amd64.exe $(SERVER_PATH)
	@echo "Multi-platform build complete in ./dist/"

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):latest .

# Docker run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env $(BINARY_NAME):latest

# Help
help:
	@echo "Available targets:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Build and run the application"
	@echo "  make dev         - Run in development mode with auto-reload"
	@echo "  make test        - Run Go unit tests"
	@echo "  make test-api    - Run API integration tests (requires USERNAME)"
	@echo "  make clean       - Remove build artifacts"
	@echo "  make fmt         - Format code"
	@echo "  make lint        - Run linter"
	@echo "  make tidy        - Tidy dependencies"
	@echo "  make deps        - Download dependencies"
	@echo "  make setup       - Setup environment (create .env)"
	@echo "  make serve       - Build and run in production mode"
	@echo "  make build-all   - Build for multiple platforms"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run  - Run Docker container"
	@echo "  make help        - Show this help message"