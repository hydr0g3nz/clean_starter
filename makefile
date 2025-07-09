.PHONY: example run sqlc-gen migration-up migration-down

example:
	@echo "Hello World!"

# Generate SQLC code
sqlc-gen:
	sqlc generate

# Run the application
run:
	go run ./cmd/main.go

# Install dependencies
install-deps:
	go mod tidy
	go mod download

# Install tools
install-tools:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Run migrations up (if you use migration tool like goose)
migration-up:
	goose -dir internal/adapter/sqlc/schema postgres "postgres://user:pass@localhost:5432/jobs?sslmode=disable" up

# Run migrations down (if you use migration tool like goose)
migration-down:
	goose -dir internal/adapter/sqlc/schema postgres "postgres://user:pass@localhost:5432/jobs?sslmode=disable" down

# Clean generated files
clean:
	rm -rf internal/adapter/sqlc/generated/*

# Build the application
build:
	go build -o bin/main ./cmd/main.go

# Run tests
test:
	go test -v ./...

# Run with hot reload (requires air)
dev:
	air

# Docker commands
docker-build:
	docker build -t topup-wallet .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

# Database setup
db-setup: migration-up sqlc-gen