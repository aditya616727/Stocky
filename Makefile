# Makefile for Stocky

.PHONY: help install build run test clean docker-up docker-down

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

install: ## Install dependencies
	go mod download

build: ## Build the application
	go build -o bin/stocky main.go

run: ## Run the application
	go run main.go

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/

docker-up: ## Start PostgreSQL with Docker
	docker run --name stocky-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=assignment -p 5432:5432 -d postgres:15

docker-down: ## Stop PostgreSQL Docker container
	docker stop stocky-postgres
	docker rm stocky-postgres

setup-db: ## Create PostgreSQL database (requires psql)
	psql -U postgres -c "CREATE DATABASE assignment;"

.DEFAULT_GOAL := help
