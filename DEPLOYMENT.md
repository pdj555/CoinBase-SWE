# Deployment Guide

## Environment Configuration

### Required Environment Variables

```bash
# JWT token signing secret (minimum 32 characters)
JWT_SECRET="your-super-secret-key-change-this-in-production"

# Token expiration time in seconds (default: 900 = 15 minutes)
TOKEN_TTL_SECONDS="900"

# HTTP server listen address (default: :8080)
HTTP_ADDR=":8080"
```

### Security Considerations

- **JWT_SECRET**: Use a cryptographically strong random string
- **HTTPS**: Always use HTTPS in production
- **Token TTL**: Balance security (shorter) vs UX (longer)

## Local Development

### Prerequisites

- Go 1.22 or later
- Make (optional, for convenience commands)

### Quick Start

```bash
# Clone repository
git clone https://github.com/coinbase/identity-service.git
cd identity-service

# Set environment variables
export JWT_SECRET="development-secret-key-12345678901234567890"
export TOKEN_TTL_SECONDS="900"

# Run the service
make run
# OR
go run ./cmd/server
```

### Development Commands

```bash
make ci        # Run all tests and static analysis
make test      # Run unit tests only
make vet       # Run static analysis
make lint      # Run linter (requires golangci-lint)
```

## Docker Deployment

### Build Image

```bash
docker build -t identity-service .
```

### Run Container

```bash
docker run -p 8080:8080 \
  -e JWT_SECRET="your-production-secret-key" \
  -e TOKEN_TTL_SECONDS="900" \
  identity-service
```

### Docker Compose

```bash
# Copy environment template
cp .env.example .env

# Edit .env with your values
vim .env

# Start service
docker-compose up -d

# Check health
curl http://localhost:8080/health

# View logs
docker-compose logs -f identity-service
```

## Production Deployment

### Container Orchestration

**Kubernetes Example**:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: identity-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: identity-service
  template:
    metadata:
      labels:
        app: identity-service
    spec:
      containers:
      - name: identity-service
        image: identity-service:latest
        ports:
        - containerPort: 8080
        env:
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: identity-secrets
              key: jwt-secret
        - name: TOKEN_TTL_SECONDS
          value: "900"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
```

### Load Balancer Configuration

**Health Check Setup**:

- **Health Check Path**: `/health`
- **Ready Check Path**: `/ready`
- **Check Interval**: 30 seconds
- **Timeout**: 10 seconds
- **Healthy Threshold**: 2 consecutive successes

### Database Migration

Currently uses in-memory storage. To add database:

1. **Implement database store**:

```go
// internal/store/postgres/user_store.go
type PostgresUserStore struct {
    db *sql.DB
}
```

1. **Update main.go**:

```go
// Replace memory store with database store
userStore := postgres.NewUserStore(db)
```

1. **Add migration system** for schema changes

## Monitoring

### Health Endpoints

**Service Health** (`/health`):

- Returns service status, uptime, version
- Use for service monitoring and alerting

**Readiness Check** (`/ready`):

- Returns when service can accept traffic
- Use for load balancer health checks

### Logging

The service provides structured JSON logs:

```json
{
  "timestamp": "2025-01-15T10:30:00Z",
  "level": "info",
  "method": "POST",
  "path": "/signin",
  "status": 200,
  "duration": "45ms",
  "ip": "192.168.1.100",
  "user_agent": "curl/7.68.0"
}
```

**Log Aggregation**: Configure log shipping to your centralized logging system (ELK, Splunk, etc.)

### Metrics

Consider adding metrics for:

- Request rate and latency
- Authentication success/failure rates
- Token generation rate
- Error rates by endpoint

**Example with Prometheus**:

```go
// Add to middleware
var (
    requestDuration = prometheus.NewHistogramVec(...)
    requestsTotal = prometheus.NewCounterVec(...)
)
```

## Security Hardening

### Environment Security

- Use secrets management (Kubernetes Secrets, AWS Secrets Manager)
- Rotate JWT secrets regularly
- Use least-privilege container permissions

### Network Security

- Deploy behind HTTPS termination
- Configure WAF rules
- Implement rate limiting at load balancer level

### Application Security

- Keep dependencies updated
- Run security scans in CI/CD
- Monitor for vulnerabilities

## Performance Tuning

### Scaling Guidelines

- **CPU**: Service is CPU-light, bcrypt is the main CPU consumer
- **Memory**: Minimal memory usage with in-memory store
- **Concurrency**: Stateless design supports horizontal scaling

### Performance Targets

- **Registration**: < 100ms (excluding bcrypt cost)
- **Login**: < 100ms (excluding bcrypt verification)
- **Token validation**: < 10ms
- **Health checks**: < 5ms

### Optimization Options

- Adjust bcrypt cost factor for security/performance balance
- Add caching layer for user lookups (with database)
- Use connection pooling for database connections

## Backup and Recovery

### Current State (In-Memory)

- No persistence, data lost on restart
- Suitable for development/testing only

### With Database

- Regular database backups
- Point-in-time recovery capability
- Test backup restoration procedures

## Troubleshooting

### Common Issues

**Service won't start**:

- Check JWT_SECRET is set and > 32 characters
- Verify port is not already in use
- Check environment variables format

**Authentication failures**:

- Verify JWT_SECRET matches between instances
- Check token expiration settings
- Validate request format and headers

**Health check failures**:

- Ensure service is fully started
- Check /health and /ready endpoints manually
- Verify load balancer configuration

### Debug Mode

```bash
# Run with verbose logging
go run ./cmd/server -v

# Check specific functionality
curl -v http://localhost:8080/health
```

### Log Analysis

```bash
# Filter authentication errors
docker-compose logs identity-service | grep "status=401"

# Monitor request performance
docker-compose logs identity-service | grep "duration"
```

## Rollback Procedure

1. **Immediate**: Route traffic to previous version
2. **Database**: Restore from backup if schema changed
3. **Configuration**: Revert environment variables if needed
4. **Verification**: Test authentication flows work correctly

## Support

For issues and questions:

- Check logs for error details
- Verify environment configuration
- Test with curl commands from API documentation
- Review deployment configuration
