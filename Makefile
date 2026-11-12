.PHONY: help build run test clean docker-up docker-down deps fmt lint migrate pm2-start pm2-restart pm2-stop pm2-logs pm2-status pm2-delete

help:
	@echo "Bali Guide Map API - Makefile Commands"
	@echo "========================================"
	@echo "make build           - Build the application"
	@echo "make run             - Run the application"
	@echo "make test            - Run tests"
	@echo "make clean           - Clean build artifacts"
	@echo "make deps            - Download dependencies"
	@echo "make fmt             - Format code"
	@echo "make lint            - Run linter"
	@echo "make docker-up       - Start Docker services"
	@echo "make docker-down     - Stop Docker services"
	@echo "make migrate         - Run database migrations"
	@echo "make pm2-start       - Build and run the API under PM2"
	@echo "make pm2-restart     - Rebuild and restart the API under PM2"
	@echo "make pm2-stop        - Stop the API under PM2"
	@echo "make pm2-logs        - Follow the API logs"
	@echo "make pm2-status      - Show PM2 process status"
	@echo "make pm2-delete      - Remove the API from PM2"

build:
	@echo "Building application..."
	@go build -o bin/api cmd/main.go

run:
	@echo "Running application..."
	@go run cmd/main.go

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

fmt:
	@echo "Formatting code..."
	@go fmt ./...

lint:
	@echo "Running linter..."
	@go vet ./...

docker-up:
	@echo "Starting Docker services..."
	@docker-compose up -d
	@echo "MySQL running on localhost:3306"
	@echo "PhpMyAdmin available at http://localhost:8081"

docker-down:
	@echo "Stopping Docker services..."
	@docker-compose down

migrate:
	@echo "Running migrations..."
	@go run cmd/main.go

# --- PM2 -------------------------------------------------------------------
# Requires PM2: npm install -g pm2

pm2-start: build
	@mkdir -p logs
	@pm2 startOrReload ecosystem.config.js
	@pm2 save

pm2-restart: build
	@pm2 restart baliguide-api --update-env

pm2-stop:
	@pm2 stop baliguide-api

pm2-logs:
	@pm2 logs baliguide-api

pm2-status:
	@pm2 status baliguide-api

pm2-delete:
	@pm2 delete baliguide-api
	@pm2 save

.DEFAULT_GOAL := help
