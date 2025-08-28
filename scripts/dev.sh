#!/bin/bash

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Check if .env exists
if [ ! -f .env ]; then
    log_info "Creating .env from example..."
    cp .env.example .env
fi

# Load environment variables
source .env

# Start database
log_info "Starting database..."
docker-compose up -d db

# Wait for database
log_info "Waiting for database..."
sleep 5

# Start API
log_info "Starting API..."
cd api && MIGRATE_ON_START=true go run cmd/api/main.go &
API_PID=$!
cd ..

# Wait for API to start
log_info "Waiting for API to start..."
sleep 3

# Start frontend
log_info "Starting frontend..."
cd web && npm run dev &
WEB_PID=$!
cd ..

log_success "Development environment started!"
echo ""
echo -e "${GREEN}🚀 Bailanysta development environment:${NC}"
echo -e "  Frontend: ${BLUE}http://localhost:5173${NC}"
echo -e "  API:      ${BLUE}http://localhost:8080${NC}"
echo ""

# Handle shutdown
cleanup() {
    echo ""
    log_info "Shutting down..."
    kill $API_PID 2>/dev/null || true
    kill $WEB_PID 2>/dev/null || true
    docker-compose stop db
    exit 0
}

trap cleanup SIGINT SIGTERM

# Wait for processes
wait
