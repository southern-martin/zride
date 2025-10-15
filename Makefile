SHELL := /bin/bash
.PHONY: help build run test clean docker-up docker-down migrate seed

# Default target
help:
	@echo "Available targets:"
	@echo "  build       - Build all services"
	@echo "  run         - Run all services locally"
	@echo "  test        - Run tests for all services"
	@echo "  clean       - Clean build artifacts"
	@echo "  docker-up   - Start all services with Docker Compose"
	@echo "  docker-down - Stop all Docker services"
	@echo "  migrate     - Run database migrations"
	@echo "  seed        - Run database seeders"
	@echo "  dev         - Start development environment"
	@echo "  lint        - Run linting for all services"
	@echo "  format      - Format code for all services"

# Load environment variables
include .env
export

# Development environment setup
dev: docker-up
	@echo "Starting development environment..."
	@echo "Services will be available at:"
	@echo "  API Gateway: http://localhost:8080"
	@echo "  Auth Service: http://localhost:8001"
	@echo "  User Service: http://localhost:8002"
	@echo "  Trip Service: http://localhost:8003"
	@echo "  Matching Service: http://localhost:8004"
	@echo "  Payment Service: http://localhost:8005"

# Build all services
build:
	@echo "Building all services..."
	cd backend/api-gateway && go build -o bin/api-gateway ./cmd/main.go
	cd backend/services/auth-service && go build -o bin/auth-service ./cmd/main.go
	cd backend/services/user-service && go build -o bin/user-service ./cmd/main.go
	cd backend/services/trip-service && go build -o bin/trip-service ./cmd/main.go
	cd backend/services/payment-service && go build -o bin/payment-service ./cmd/main.go
	cd backend/services/matching-service && pip install -r requirements.txt

# Run services locally (without Docker)
run:
	@echo "Starting services locally..."
	@trap 'kill %1 %2 %3 %4 %5 %6; exit' INT; \
	cd backend/api-gateway && go run cmd/main.go & \
	cd backend/services/auth-service && go run cmd/main.go & \
	cd backend/services/user-service && go run cmd/main.go & \
	cd backend/services/trip-service && go run cmd/main.go & \
	cd backend/services/payment-service && go run cmd/main.go & \
	cd backend/services/matching-service && python app.py & \
	wait

# Run tests
test:
	@echo "Running tests..."
	cd backend/api-gateway && go test ./...
	cd backend/services/auth-service && go test ./...
	cd backend/services/user-service && go test ./...
	cd backend/services/trip-service && go test ./...
	cd backend/services/payment-service && go test ./...
	cd backend/services/matching-service && python -m pytest tests/

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	find . -name "bin" -type d -exec rm -rf {} +
	find . -name "__pycache__" -type d -exec rm -rf {} +
	find . -name "*.pyc" -delete
	find . -name ".coverage" -delete
	docker system prune -f

# Docker operations
docker-up:
	@echo "Starting Docker services..."
	docker-compose up -d
	@echo "Waiting for services to be healthy..."
	@sleep 10
	@echo "Services started successfully!"

docker-down:
	@echo "Stopping Docker services..."
	docker-compose down

docker-rebuild:
	@echo "Rebuilding Docker services..."
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

# Database operations
migrate:
	@echo "Running database migrations..."
	cd backend/services/auth-service && go run cmd/migrate/main.go up
	cd backend/services/user-service && go run cmd/migrate/main.go up
	cd backend/services/trip-service && go run cmd/migrate/main.go up
	cd backend/services/payment-service && go run cmd/migrate/main.go up

migrate-down:
	@echo "Rolling back database migrations..."
	cd backend/services/payment-service && go run cmd/migrate/main.go down
	cd backend/services/trip-service && go run cmd/migrate/main.go down
	cd backend/services/user-service && go run cmd/migrate/main.go down
	cd backend/services/auth-service && go run cmd/migrate/main.go down

seed:
	@echo "Running database seeders..."
	cd database/seeds && go run seed.go

# Code quality
lint:
	@echo "Running linting..."
	cd backend/api-gateway && golangci-lint run
	cd backend/services/auth-service && golangci-lint run
	cd backend/services/user-service && golangci-lint run
	cd backend/services/trip-service && golangci-lint run
	cd backend/services/payment-service && golangci-lint run
	cd backend/services/matching-service && flake8 .

format:
	@echo "Formatting code..."
	cd backend/api-gateway && go fmt ./...
	cd backend/services/auth-service && go fmt ./...
	cd backend/services/user-service && go fmt ./...
	cd backend/services/trip-service && go fmt ./...
	cd backend/services/payment-service && go fmt ./...
	cd backend/services/matching-service && black . && isort .

# Frontend operations
frontend-install:
	@echo "Installing frontend dependencies..."
	cd frontend/zalo-mini-app && npm install

frontend-dev:
	@echo "Starting frontend development server..."
	cd frontend/zalo-mini-app && npm start

frontend-build:
	@echo "Building frontend for production..."
	cd frontend/zalo-mini-app && npm run build

# Logs
logs:
	@echo "Showing logs for all services..."
	docker-compose logs -f

logs-service:
	@echo "Showing logs for $(SERVICE)..."
	docker-compose logs -f $(SERVICE)

# Health checks
health:
	@echo "Checking service health..."
	@curl -s http://localhost:8080/health || echo "API Gateway: DOWN"
	@curl -s http://localhost:8001/health || echo "Auth Service: DOWN"
	@curl -s http://localhost:8002/health || echo "User Service: DOWN"
	@curl -s http://localhost:8003/health || echo "Trip Service: DOWN"
	@curl -s http://localhost:8004/health || echo "Matching Service: DOWN"
	@curl -s http://localhost:8005/health || echo "Payment Service: DOWN"

# Documentation
docs:
	@echo "Generating API documentation..."
	cd backend/api-gateway && swag init
	cd backend/services/auth-service && swag init
	cd backend/services/user-service && swag init
	cd backend/services/trip-service && swag init
	cd backend/services/payment-service && swag init

# Setup development environment
setup:
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then cp .env.example .env; echo "Created .env file from example"; fi
	@echo "Please update .env file with your configuration"
	make frontend-install
	@echo "Development environment setup complete!"
	@echo "Run 'make dev' to start the development server"