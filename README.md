# Identity Service

**Production-ready user authentication microservice built for scale.**

A secure, fast, and maintainable identity service that handles user registration, authentication, and session management. Built with Go 1.22, designed with Coinbase-grade security standards.

## What it does

- **User Registration**: Secure signup with email validation and password strength requirements
- **Authentication**: JWT-based login with bcrypt password hashing  
- **Session Management**: Protected routes with token validation
- **Health Monitoring**: Built-in health checks for production deployments
- **Input Validation**: Comprehensive request validation with clear error messages

## Quick Start

**Prerequisites**: Go 1.22+

```bash
# Clone and setup
git clone https://github.com/coinbase/identity-service.git
cd identity-service

# Set environment variables
export JWT_SECRET="your-super-secret-key-min-32-chars"
export TOKEN_TTL_SECONDS="900"

# Run the service
make run
```

The service starts on `http://localhost:8080`

## API Reference

### Public Endpoints

#### Register User

```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepass123"}'
```

#### Login

```bash
curl -X POST http://localhost:8080/signin \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepass123"}'
```

### Protected Endpoints

#### Get User Profile
(requires Authorization header)

```bash
curl -X GET http://localhost:8080/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Health Endpoints

#### Service Health

```bash
curl http://localhost:8080/health
# Returns: {"status":"healthy","uptime":"2m30s",...}
```

#### Readiness Check

```bash
curl http://localhost:8080/ready
# Returns: {"status":"ready"}
```

## Architecture

```markdown
┌─────────────────┐
│   HTTP Router   │  ← Handles requests, middleware
└─────────┬───────┘
          │
┌─────────▼───────┐
│ Auth Handlers   │  ← Request validation, response formatting  
└─────────┬───────┘
          │
┌─────────▼───────┐
│ Auth Service    │  ← Business logic, user operations
└─────────┬───────┘
          │
┌─────────▼───────┐
│   Data Store    │  ← User persistence (in-memory/database)
└─────────────────┘
```

#### Key Design Principles

- **Dependency Injection**: Easy to test and extend
- **Interface-Based**: Swap implementations without code changes
- **Clear Separation**: HTTP, business logic, and data layers are distinct
- **Security First**: Input validation, secure defaults, proper error handling

## Security Features

- **Password Security**: bcrypt hashing with salt
- **Input Validation**: Email format validation, password strength requirements
- **JWT Tokens**: Secure token generation with configurable expiration
- **Request Validation**: Comprehensive input sanitization
- **Error Handling**: Secure error messages without information leakage

## Development

#### Run Tests

```bash
make ci                    # Run all tests + static analysis
go test -v ./...          # Run tests with verbose output
go test -cover ./...      # Run tests with coverage report
```

#### Project Structure

```markdown
├── cmd/server/            # Application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── handler/          # HTTP request handlers
│   ├── middleware/       # HTTP middleware (logging, auth)
│   ├── model/            # Data models
│   ├── server/           # HTTP server setup
│   ├── service/          # Business logic
│   ├── store/            # Data persistence layer
│   └── validator/        # Input validation
├── pkg/
│   ├── hash/             # Password hashing utilities
│   └── token/            # JWT token management
└── docker-compose.yml    # Container deployment
```

## Production Deployment

#### Using Docker Compose

```bash
cp .env.example .env      # Configure environment variables
docker-compose up -d      # Start service with health checks
```

#### Environment Variables

```bash
HTTP_ADDR=:8080                    # Server listen address
JWT_SECRET=your-secret-key         # JWT signing secret (32+ chars)
TOKEN_TTL_SECONDS=900             # Token expiration (15 minutes)
```

## Performance

- **Fast**: Sub-100ms response times for authentication operations
- **Concurrent**: Handles multiple simultaneous requests safely
- **Memory Efficient**: Minimal memory footprint with efficient data structures
- **Observable**: Structured logging with request duration tracking

## Testing

The service includes comprehensive testing:

- **28 Unit Tests**: 57% overall code coverage
- **Integration Tests**: Full API endpoint testing
- **Security Tests**: Input validation and edge case handling
- **Performance Tests**: Concurrent request handling

#### Test Coverage by Component

- Hash utilities: 100%
- JWT tokens: 100%  
- Data store: 100%
- Input validation: 100%
- Auth service: 87.5%
- HTTP handlers: 88.9%

## License

MIT License - see [LICENSE](LICENSE) for details.
