# note: call scripts from /scripts

.PHONY: help build run test test-unit test-integration clean docker-up docker-down swagger

help: ## Display this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building application..."
	go build -o bin/todoapp cmd/todoapp/main.go

run: ## Run the application
	@echo "Running application..."
	go run cmd/todoapp/main.go

test: test-unit test-integration ## Run all tests

test-unit: ## Run unit tests
	@echo "Running unit tests..."
	go test -mod=mod ./internal/app/todoapp/services/... -v -cover

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	go test -mod=mod ./test/integration/... -v -cover

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -f todoapp

docker-up: ## Start Docker containers
	@echo "Starting Docker containers..."
	docker-compose up -d

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger docs..."
	swag init -g cmd/todoapp/main.go -o api/docs

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .
