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

# Check if docker and docker-compose are installed
check_dependencies() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed. Please install Docker first."
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        log_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
}

# Load environment variables
load_env() {
    if [ -f ".env" ]; then
        log_info "Loading environment variables from .env"
        export $(grep -v '^#' .env | xargs)
    else
        log_warning "No .env file found. Using default values."
    fi
}

# Build and start services
deploy_dev() {
    log_info "Starting development environment..."
    docker-compose up -d --build
    log_success "Development environment started!"
    log_info "Frontend: http://localhost:3000"
    log_info "Backend API: http://localhost:8080"
    log_info "Database: localhost:5432"
}

# Deploy production environment
deploy_prod() {
    if [ ! -f ".env" ]; then
        log_error "Production deployment requires .env file with configuration"
        log_info "Copy env.example to .env and configure your settings"
        exit 1
    fi

    log_info "Starting production environment..."
    docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
    log_success "Production environment started!"
    log_info "Application: http://localhost"
    log_info "Check logs with: docker-compose logs -f"
}

# Stop services
stop_services() {
    log_info "Stopping all services..."
    docker-compose down
    log_success "Services stopped!"
}

# View logs
show_logs() {
    if [ -z "$2" ]; then
        docker-compose logs -f
    else
        docker-compose logs -f $2
    fi
}

# Build services
build_services() {
    log_info "Building services..."
    docker-compose build --no-cache
    log_success "Services built!"
}

# Run database migrations
run_migrations() {
    log_info "Running database migrations..."
    docker-compose exec api ./bailanysta-api --migrate-only
    log_success "Migrations completed!"
}

# Backup database
backup_db() {
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local backup_file="backup_${timestamp}.sql"

    log_info "Creating database backup..."
    docker-compose exec db pg_dump -U bailanysta_user bailanysta > "infra/backups/${backup_file}"
    log_success "Backup created: infra/backups/${backup_file}"
}

# Restore database
restore_db() {
    if [ -z "$2" ]; then
        log_error "Please provide backup file path: ./deploy.sh restore /path/to/backup.sql"
        exit 1
    fi

    local backup_file="$2"
    if [ ! -f "$backup_file" ]; then
        log_error "Backup file not found: $backup_file"
        exit 1
    fi

    log_warning "This will replace the current database. Are you sure? (y/N)"
    read -r confirm
    if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
        log_info "Database restore cancelled."
        exit 0
    fi

    log_info "Restoring database from $backup_file..."
    docker-compose exec -T db psql -U bailanysta_user bailanysta < "$backup_file"
    log_success "Database restored!"
}

# Clean up
cleanup() {
    log_info "Cleaning up Docker resources..."
    docker-compose down -v --remove-orphans
    docker system prune -f
    log_success "Cleanup completed!"
}

# Show help
show_help() {
    echo "Bailanysta Deployment Script"
    echo ""
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  dev       Start development environment"
    echo "  prod      Start production environment"
    echo "  stop      Stop all services"
    echo "  logs      Show logs (use 'logs <service>' for specific service)"
    echo "  build     Build all services"
    echo "  migrate   Run database migrations"
    echo "  backup    Create database backup"
    echo "  restore   Restore database from backup file"
    echo "  cleanup   Clean up Docker resources"
    echo "  help      Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 dev                    # Start development environment"
    echo "  $0 prod                   # Start production environment"
    echo "  $0 logs api               # Show API logs"
    echo "  $0 backup                 # Create database backup"
    echo "  $0 restore backup.sql     # Restore database from backup"
}

# Main script logic
main() {
    check_dependencies
    load_env

    case "${1:-help}" in
        "dev")
            deploy_dev
            ;;
        "prod")
            deploy_prod
            ;;
        "stop")
            stop_services
            ;;
        "logs")
            show_logs "$@"
            ;;
        "build")
            build_services
            ;;
        "migrate")
            run_migrations
            ;;
        "backup")
            mkdir -p infra/backups
            backup_db
            ;;
        "restore")
            restore_db "$@"
            ;;
        "cleanup")
            cleanup
            ;;
        "help"|*)
            show_help
            ;;
    esac
}

# Run main function with all arguments
main "$@"
