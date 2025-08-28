#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if .env file exists
check_env_file() {
    if [ ! -f .env ]; then
        log_warning ".env file not found. Creating from example..."
        if [ -f .env.example ]; then
            cp .env.example .env
            log_info "Please edit .env file with your configuration before continuing."
            echo "Press Enter to continue after editing .env, or Ctrl+C to cancel..."
            read -r
        else
            log_error ".env.example not found. Cannot create .env file."
            exit 1
        fi
    fi
}

# Check if docker and docker-compose are available
check_dependencies() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed. Please install Docker first."
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        log_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi

    log_success "Docker and Docker Compose are available"
}

# Wait for database to be ready
wait_for_db() {
    log_info "Waiting for database to be ready..."
    
    local max_attempts=30
    local attempt=0
    
    while [ $attempt -lt $max_attempts ]; do
        if docker-compose exec -T db pg_isready -U bailanysta_user -d bailanysta > /dev/null 2>&1; then
            log_success "Database is ready!"
            return 0
        fi
        
        attempt=$((attempt + 1))
        log_info "Waiting for database... (attempt $attempt/$max_attempts)"
        sleep 2
    done
    
    log_error "Database failed to start within expected time"
    return 1
}

# Run migrations
run_migrations() {
    log_info "Running database migrations..."
    
    # Set environment variable to enable migrations
    export MIGRATE_ON_START=true
    
    # Restart API service to run migrations
    docker-compose restart api
    
    # Wait a bit for migrations to complete
    sleep 5
    
    log_success "Migrations completed"
}

# Health check
health_check() {
    log_info "Performing health checks..."
    
    # Check API health
    if curl -s -f http://localhost:8080/health > /dev/null; then
        log_success "API is healthy"
    else
        log_warning "API health check failed"
    fi
    
    # Check frontend
    if curl -s -f http://localhost:3000/ > /dev/null; then
        log_success "Frontend is accessible"
    else
        log_warning "Frontend health check failed"
    fi
}

# Main deployment function
deploy() {
    local mode=${1:-"dev"}
    
    log_info "Starting Bailanysta deployment (mode: $mode)..."
    
    # Pre-flight checks
    check_dependencies
    check_env_file
    
    # Stop any existing services
    log_info "Stopping existing services..."
    docker-compose down
    
    # Build and start services
    log_info "Building and starting services..."
    if [ "$mode" = "prod" ]; then
        docker-compose up --build -d
    else
        docker-compose up --build -d
    fi
    
    # Wait for database
    wait_for_db
    
    # Run migrations
    run_migrations
    
    # Perform health checks
    sleep 10
    health_check
    
    log_success "Deployment completed successfully!"
    echo ""
    echo -e "${GREEN}🚀 Bailanysta is now running:${NC}"
    echo -e "  Frontend: ${BLUE}http://localhost:3000${NC}"
    echo -e "  API:      ${BLUE}http://localhost:8080${NC}"
    echo -e "  Health:   ${BLUE}http://localhost:8080/health${NC}"
    echo ""
    echo -e "${YELLOW}To view logs:${NC} docker-compose logs -f"
    echo -e "${YELLOW}To stop:${NC}     docker-compose down"
}

# Script usage
usage() {
    echo "Usage: $0 [dev|prod]"
    echo ""
    echo "Commands:"
    echo "  dev   - Deploy in development mode (default)"
    echo "  prod  - Deploy in production mode"
    echo ""
    echo "Examples:"
    echo "  $0        # Deploy in development mode"
    echo "  $0 dev    # Deploy in development mode"
    echo "  $0 prod   # Deploy in production mode"
}

# Main
case "${1:-dev}" in
    "dev"|"prod")
        deploy "$1"
        ;;
    "help"|"-h"|"--help")
        usage
        ;;
    *)
        log_error "Unknown command: $1"
        usage
        exit 1
        ;;
esac
