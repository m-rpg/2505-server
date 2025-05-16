.PHONY: build run test clean migrate deps fmt lint docs help

# Build the application
build:
	go build -o bin/server

# Run the application
run:
	go run main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Run database migrations
migrate:
	go run main.go migrate

# Install dependencies
deps:
	go mod download

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Generate API documentation
docs:
	swag init -g main.go

# Help command
help:
	@echo "Available commands:"
	@echo "  make build    - Build the application"
	@echo "  make run      - Run the application"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make migrate  - Run database migrations"
	@echo "  make deps     - Install dependencies"
	@echo "  make fmt      - Format code"
	@echo "  make lint     - Run linter"
	@echo "  make docs     - Generate API documentation"
