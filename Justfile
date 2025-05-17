# Build the application
build:
    go build -o bin/server

# Run the application
run *ARGS:
    go run main.go {{ARGS}}

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
