.PHONY: help build test clean dev prod stop logs migrate backup

# Default target
help: ## Show this help message
	@echo "Bailanysta - Educational Platform"
	@echo ""
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development commands
dev: ## Start development environment
	./infra/docker/deploy.sh dev

prod: ## Start production environment
	./infra/docker/deploy.sh prod

stop: ## Stop all services
	./infra/docker/deploy.sh stop

logs: ## Show all logs
	./infra/docker/deploy.sh logs

logs-%: ## Show logs for specific service (e.g., make logs-api)
	./infra/docker/deploy.sh logs $*

build: ## Build all services
	./infra/docker/deploy.sh build

# Database commands
migrate: ## Run database migrations
	./infra/docker/deploy.sh migrate

backup: ## Create database backup
	./infra/docker/deploy.sh backup

# Testing commands
test: ## Run all tests
	go test ./api/... -v
	cd web && npm test -- --run

test-go: ## Run Go tests only
	go test ./api/... -v

test-go-%: ## Run specific Go package tests
	go test ./api/internal/$*/... -v

test-frontend: ## Run frontend tests only
	cd web && npm test -- --run

test-frontend-watch: ## Run frontend tests in watch mode
	cd web && npm test

# Code quality
lint: ## Run linters
	cd web && npm run lint
	# Add Go linter here if needed

format: ## Format code
	cd web && npm run format
	# Add Go formatter here if needed

# Cleanup
clean: ## Clean up Docker resources
	./infra/docker/deploy.sh cleanup

clean-all: ## Clean up everything including volumes
	docker-compose down -v --remove-orphans
	docker system prune -f --volumes

# Installation and setup
install: ## Install all dependencies
	go mod tidy
	cd web && npm install

setup: ## Initial project setup
	@echo "Setting up Bailanysta project..."
	make install
	@echo "Setup complete! Run 'make dev' to start development environment."

# Docker utilities
docker-clean: ## Remove all Docker containers, images, and volumes
	docker stop $$(docker ps -aq) 2>/dev/null || true
	docker rm $$(docker ps -aq) 2>/dev/null || true
	docker rmi $$(docker images -q) 2>/dev/null || true
	docker volume rm $$(docker volume ls -q) 2>/dev/null || true
	docker system prune -f --volumes

# Health checks
health: ## Check health of all services
	@echo "Checking API health..."
	curl -f http://localhost:8080/health || echo "API is not healthy"
	@echo "Checking frontend..."
	curl -f http://localhost:3000/ || echo "Frontend is not healthy"

# Production deployment helpers
deploy-check: ## Pre-deployment checks
	@echo "Running pre-deployment checks..."
	@test -f .env || (echo "Error: .env file not found" && exit 1)
	@grep -q "JWT_SECRET" .env || (echo "Error: JWT_SECRET not set in .env" && exit 1)
	@grep -q "DB_PASSWORD" .env || (echo "Error: DB_PASSWORD not set in .env" && exit 1)
	@echo "Pre-deployment checks passed!"

# Development helpers
seed: ## Seed database with test data
	@echo "Seeding database..."
	docker-compose exec db psql -U bailanysta_user -d bailanysta -f /docker-entrypoint-initdb.d/init.sql

reset-db: ## Reset database (WARNING: This will delete all data)
	@echo "WARNING: This will delete all database data!"
	@read -p "Are you sure? (y/N) " confirm && [ "$$confirm" = "y" ] || exit 1
	docker-compose exec db psql -U bailanysta_user -d bailanysta -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

# Git hooks
install-hooks: ## Install git hooks
	@echo "Installing git hooks..."
	@mkdir -p .git/hooks
	@cp infra/git-hooks/pre-commit .git/hooks/pre-commit 2>/dev/null || echo "No pre-commit hook found"
	@chmod +x .git/hooks/*
	@echo "Git hooks installed!"
