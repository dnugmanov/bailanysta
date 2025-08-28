.PHONY: help dev stop restart logs clean health build

# Default target
help: ## Show this help message
	@echo "🚀 Bailanysta - Educational Platform"
	@echo ""
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development commands
dev: ## Start development environment
	@echo "🔧 Starting development environment..."
	docker-compose up --build -d

prod: ## Start production environment with SSL
	@echo "🚀 Starting production environment..."
	@if [ ! -f .env ]; then echo "❌ Error: .env file not found. Copy .env.example to .env and configure it."; exit 1; fi
	@if [ -z "$$SSL_CERT_PATH" ] || [ -z "$$SSL_KEY_PATH" ]; then echo "❌ Error: SSL_CERT_PATH and SSL_KEY_PATH must be set in .env"; exit 1; fi
	docker-compose -f docker-compose.prod.yml up --build -d

stop: ## Stop all services
	@echo "🛑 Stopping all services..."
	docker-compose down

stop-prod: ## Stop production services
	@echo "🛑 Stopping production services..."
	docker-compose -f docker-compose.prod.yml down

restart: ## Restart all services
	@echo "🔄 Restarting all services..."
	docker-compose restart

# Logs
logs: ## Show logs for all services
	docker-compose logs -f

logs-api: ## Show API logs
	docker-compose logs -f api

logs-web: ## Show web logs
	docker-compose logs -f web

logs-db: ## Show database logs
	docker-compose logs -f db

logs-prod: ## Show production logs
	docker-compose -f docker-compose.prod.yml logs -f

logs-nginx: ## Show nginx logs (production)
	docker-compose -f docker-compose.prod.yml logs -f nginx

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

# Health check
health: ## Check health of all services
	@echo "🩺 Checking services health..."
	@echo "API Health:"
	@curl -s http://localhost:8080/health | jq . || echo "API not responding"
	@echo "Frontend:"
	@curl -s -o /dev/null -w "Status: %{http_code}" http://localhost:3000/ || echo "Frontend not responding"
	@echo ""

# Cleanup
clean: ## Clean up containers and images
	@echo "🧹 Cleaning up..."
	docker-compose down --rmi all --volumes --remove-orphans

clean-prod: ## Clean up production containers and images
	@echo "🧹 Cleaning up production..."
	docker-compose -f docker-compose.prod.yml down --rmi all --volumes --remove-orphans