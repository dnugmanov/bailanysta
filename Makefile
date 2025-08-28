.PHONY: help install build start stop restart logs clean dev prod migrate

# Default target
help: ## Show this help message
	@echo "🚀 Bailanysta - Educational Platform"
	@echo ""
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Installation
install: ## Install dependencies
	@echo "📦 Installing dependencies..."
	go mod tidy
	cd web && npm install

# Development commands
dev: ## Start development environment
	@echo "🔧 Starting development environment..."
	docker-compose up --build -d

start: ## Start all services
	@echo "🚀 Starting all services..."
	docker-compose up -d

stop: ## Stop all services
	@echo "🛑 Stopping all services..."
	docker-compose down

restart: ## Restart all services
	@echo "🔄 Restarting all services..."
	docker-compose restart

logs: ## Show logs for all services
	docker-compose logs -f

logs-api: ## Show API logs
	docker-compose logs -f api

logs-web: ## Show web logs
	docker-compose logs -f web

logs-db: ## Show database logs
	docker-compose logs -f db

# Build commands
build: ## Build all services
	@echo "🔨 Building all services..."
	docker-compose build

build-api: ## Build API service
	@echo "🔨 Building API..."
	docker-compose build api

build-web: ## Build web service
	@echo "🔨 Building web..."
	docker-compose build web

# Production deployment
prod: ## Deploy to production
	@echo "🚀 Deploying to production..."
	@if [ ! -f .env ]; then echo "❌ Error: .env file not found. Copy .env.example to .env and configure it."; exit 1; fi
	docker-compose -f docker-compose.yml up --build -d

# Database commands
migrate: ## Run database migrations
	@echo "🗃️ Running database migrations..."
	docker-compose exec api ./bailanysta-api -migrate

migrate-force: ## Force run migrations (use with caution)
	@echo "⚠️ Force running database migrations..."
	docker-compose exec db psql -U bailanysta_user -d bailanysta -f /migrations/0001_init.sql

# Utility commands
shell-api: ## Access API container shell
	docker-compose exec api sh

shell-db: ## Access database shell
	docker-compose exec db psql -U bailanysta_user -d bailanysta

shell-web: ## Access web container shell
	docker-compose exec web sh

# Health checks
health: ## Check health of all services
	@echo "🩺 Checking services health..."
	@echo "API Health:"
	@curl -s http://localhost:8080/health | jq . || echo "API not responding"
	@echo "Frontend:"
	@curl -s -o /dev/null -w "Status: %{http_code}" http://localhost:3000/ || echo "Frontend not responding"
	@echo ""

# Cleanup commands
clean: ## Clean up containers and images
	@echo "🧹 Cleaning up..."
	docker-compose down --rmi all --volumes --remove-orphans

clean-all: ## Remove everything including volumes
	@echo "🧹 Removing everything..."
	docker-compose down -v --remove-orphans
	docker system prune -f --volumes

# Testing commands
test: ## Run all tests
	@echo "🧪 Running tests..."
	go test ./api/... -v
	cd web && npm run test:run

test-api: ## Run API tests only
	go test ./api/... -v

test-web: ## Run frontend tests only
	cd web && npm run test:run

# Code quality
lint: ## Run linters
	@echo "🔍 Running linters..."
	cd web && npm run lint

format: ## Format code
	@echo "✨ Formatting code..."
	cd web && npm run format

# Quick setup for new environment
setup: ## Initial setup for new environment
	@echo "🛠️ Setting up Bailanysta..."
	@if [ ! -f .env ]; then echo "📝 Creating .env from example..."; cp .env.example .env; echo "⚠️ Please edit .env file with your configuration"; fi
	make install
	@echo "✅ Setup complete! Run 'make dev' to start development or 'make prod' for production."