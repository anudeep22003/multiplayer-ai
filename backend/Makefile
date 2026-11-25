LOCAL_DB_URL ?= postgres://multiplayer_ai:multiplayer_ai@localhost:5432/multiplayer_ai?sslmode=disable

.PHONY: help db-up db-down migrate-status migrate-up migrate-down sqlc dev test test-integration build

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-18s %s\n", $$1, $$2}'

db-up: ## Start local PostgreSQL container
	docker compose up -d --wait

db-down: ## Stop local PostgreSQL container
	docker compose down

migrate-status: ## Show goose migration status
	go run ./cmd/migrate status

migrate-up: ## Apply pending goose migrations
	go run ./cmd/migrate up

migrate-down: ## Roll back the latest goose migration
	go run ./cmd/migrate down

sqlc: ## Regenerate Go database code from SQL queries
	sqlc generate

dev: ## Run backend HTTP server locally
	go run ./cmd/server

test: ## Run go vet and unit tests with race detector
	go vet ./...
	go test -race ./...

test-integration: ## Run tests including DB integration tests against LOCAL_DB_URL
	TEST_DATABASE_URL=$(LOCAL_DB_URL) go test -race ./...

build: ## Build all backend packages
	go build ./...
