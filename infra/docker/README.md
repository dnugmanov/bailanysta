# Bailanysta Deployment Guide

This guide provides instructions for deploying the Bailanysta educational platform using Docker and Docker Compose.

## Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+
- Git
- At least 4GB RAM available

## Quick Start

### Development Environment

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd bailanysta
   ```

2. **Start development environment:**
   ```bash
   ./infra/docker/deploy.sh dev
   ```

3. **Access the application:**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Database: localhost:5432 (user: bailanysta_user, password: bailanysta_password)

### Production Environment

1. **Configure environment variables:**
   ```bash
   cp infra/docker/env.example .env
   # Edit .env with your production settings
   ```

2. **Start production environment:**
   ```bash
   ./infra/docker/deploy.sh prod
   ```

3. **Access the application:**
   - Application: http://localhost
   - Nginx will proxy requests to the appropriate services

## Available Commands

The deployment script provides several commands:

```bash
# Development
./infra/docker/deploy.sh dev       # Start development environment
./infra/docker/deploy.sh stop      # Stop all services
./infra/docker/deploy.sh logs      # Show all logs
./infra/docker/deploy.sh logs api  # Show API logs only

# Production
./infra/docker/deploy.sh prod      # Start production environment

# Maintenance
./infra/docker/deploy.sh build     # Build all services
./infra/docker/deploy.sh migrate   # Run database migrations
./infra/docker/deploy.sh backup    # Create database backup
./infra/docker/deploy.sh cleanup   # Clean up Docker resources
```

## Services Architecture

```
Internet
    ↓
[Nginx Reverse Proxy] (port 80/443)
    ↓
├── [React Frontend] (port 3000)
├── [Go Backend API] (port 8080)
├── [PostgreSQL] (port 5432)
└── [Redis] (port 6379)
```

## Environment Variables

### Required Variables

- `DB_USER` - PostgreSQL username
- `DB_PASSWORD` - PostgreSQL password
- `JWT_SECRET` - JWT signing secret (min 32 characters)
- `OPENAI_API_KEY` - OpenAI API key for AI features

### Optional Variables

- `JWT_EXPIRY` - JWT token expiry (default: 15m)
- `REFRESH_EXPIRY` - Refresh token expiry (default: 168h)
- `CORS_ORIGIN` - CORS allowed origin (default: http://localhost:3000)
- `LOG_LEVEL` - Logging level (default: info)
- `RATE_LIMIT_RPM` - API rate limit (default: 100)

## Database Configuration

The application uses PostgreSQL with the following extensions:
- `uuid-ossp` - For UUID generation
- `pg_trgm` - For full-text search

Database migrations are run automatically on startup.

## SSL/HTTPS Configuration

For production deployment with SSL:

1. Obtain SSL certificate (Let's Encrypt, etc.)
2. Place certificate files in `infra/ssl/`
3. Update nginx configuration with certificate paths
4. Uncomment HTTPS server block in nginx.conf

## Monitoring and Logging

### Logs

View logs for all services:
```bash
docker-compose logs -f
```

View logs for specific service:
```bash
docker-compose logs -f api
docker-compose logs -f web
```

### Health Checks

All services include health checks:
- Database: PostgreSQL connection check
- API: HTTP health endpoint
- Frontend: Served by Nginx

## Backup and Restore

### Database Backup

```bash
./infra/docker/deploy.sh backup
```

This creates a backup file in `infra/backups/`

### Database Restore

```bash
./infra/docker/deploy.sh restore infra/backups/backup_20231201.sql
```

**⚠️ WARNING:** This will replace the current database!

## Scaling

### Production Scaling

For high-traffic production environments:

1. **Database:** Use managed PostgreSQL (AWS RDS, Google Cloud SQL)
2. **Redis:** Use managed Redis (AWS ElastiCache, Google Memorystore)
3. **API:** Scale horizontally with load balancer
4. **Frontend:** Use CDN for static assets

### Docker Swarm/Kubernetes

For larger deployments, consider:
- Docker Swarm for simple orchestration
- Kubernetes for complex microservices
- Ingress controllers for load balancing

## Troubleshooting

### Common Issues

1. **Port conflicts:**
   - Change ports in docker-compose.yml
   - Use `docker-compose ps` to check running containers

2. **Database connection issues:**
   - Check database logs: `docker-compose logs db`
   - Verify environment variables
   - Ensure database is healthy: `docker-compose exec db pg_isready`

3. **API not responding:**
   - Check API logs: `docker-compose logs api`
   - Verify JWT_SECRET is set
   - Check database connectivity

4. **Frontend build issues:**
   - Clear node_modules: `docker-compose exec web rm -rf node_modules`
   - Rebuild: `docker-compose up --build web`

### Debug Mode

For debugging Go applications:
```bash
# Enable debugging in docker-compose.override.yml
ports:
  - "2345:2345"  # Delve debugger port
```

Then attach debugger to localhost:2345

## Security Considerations

### Production Checklist

- [ ] Change default passwords
- [ ] Use strong JWT_SECRET (min 32 characters)
- [ ] Enable SSL/HTTPS
- [ ] Configure firewall rules
- [ ] Set up monitoring and alerts
- [ ] Regular backup schedule
- [ ] Keep dependencies updated
- [ ] Use environment-specific configurations

### Secrets Management

For production, use:
- Docker secrets
- Environment variables
- External secret managers (AWS Secrets Manager, HashiCorp Vault)

## Performance Optimization

### Database

- Use connection pooling
- Enable query logging for optimization
- Set appropriate memory limits
- Use indexes for common queries

### API

- Enable gzip compression
- Set appropriate rate limits
- Use caching where appropriate
- Monitor response times

### Frontend

- Enable gzip compression in nginx
- Use CDN for static assets
- Implement code splitting
- Optimize bundle size

## Support

For issues and questions:
1. Check the logs: `docker-compose logs`
2. Review Docker documentation
3. Check GitHub issues
4. Review application documentation
