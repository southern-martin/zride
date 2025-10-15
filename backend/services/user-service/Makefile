# Variables
SERVICE_NAME=user-service
DOCKER_IMAGE=zride/$(SERVICE_NAME)
GO_VERSION=1.21

# Default target
.DEFAULT_GOAL := help

# Build the application
.PHONY: build
build:
	@echo "Building $(SERVICE_NAME)..."
	go build -o bin/$(SERVICE_NAME) .

# Run the application
.PHONY: run
run:
	@echo "Running $(SERVICE_NAME)..."
	go run .

# Test the application
.PHONY: test
test:
	@echo "Running tests for $(SERVICE_NAME)..."
	go test -v ./...

# Test with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage for $(SERVICE_NAME)..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	golangci-lint run

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Build Docker image
.PHONY: docker-build
docker-build:
	@echo "Building Docker image $(DOCKER_IMAGE)..."
	docker build -t $(DOCKER_IMAGE):latest .

# Run Docker container
.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8082:8082 --env-file .env $(DOCKER_IMAGE):latest

# Push Docker image
.PHONY: docker-push
docker-push:
	@echo "Pushing Docker image $(DOCKER_IMAGE)..."
	docker push $(DOCKER_IMAGE):latest

# Development setup
.PHONY: dev-setup
dev-setup:
	@echo "Setting up development environment..."
	go mod download
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Live reload for development
.PHONY: dev
dev:
	@echo "Starting live reload development server..."
	air

# Database migration (if needed)
.PHONY: migrate-up
migrate-up:
	@echo "Running database migrations..."
	# Add migration commands here

.PHONY: migrate-down
migrate-down:
	@echo "Rolling back database migrations..."
	# Add rollback commands here

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  clean          - Clean build artifacts"
	@echo "  fmt            - Format code"
	@echo "  lint           - Lint code"
	@echo "  deps           - Download and tidy dependencies"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  docker-push    - Push Docker image"
	@echo "  dev-setup      - Setup development environment"
	@echo "  dev            - Start live reload development"
	@echo "  migrate-up     - Run database migrations"
	@echo "  migrate-down   - Rollback database migrations"
	@echo "  help           - Show this help message"