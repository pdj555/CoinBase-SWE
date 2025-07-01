# Architecture Overview

## System Design

The Identity Service follows a **layered architecture** with clear separation of concerns, making it maintainable, testable, and extensible.

## Core Components

### 1. HTTP Layer (`internal/server`, `internal/handler`)

**Responsibility**: Handle HTTP requests, route to business logic, format responses

**Key Files**:

- `server/router.go` - HTTP routing and middleware setup
- `handler/auth.go` - Authentication endpoint handlers
- `handler/health.go` - Health check endpoints

**Design Decisions**:

- Uses Gorilla Mux for routing (industry standard, flexible)
- Middleware-based architecture for cross-cutting concerns
- Consistent JSON responses with proper HTTP status codes

### 2. Business Logic Layer (`internal/service`)

**Responsibility**: Core business rules, user operations, security logic

**Key Files**:

- `service/auth.go` - Authentication business logic

**Design Decisions**:

- Interface-based design for easy testing and mocking
- Dependency injection for external dependencies
- Clear error types for different failure scenarios

### 3. Data Layer (`internal/store`)

**Responsibility**: User data persistence and retrieval

**Key Files**:

- `store/store.go` - Data access interface
- `store/memory/user_store.go` - In-memory implementation

**Design Decisions**:

- Interface-first design enables easy database swapping
- Thread-safe operations for concurrent access
- Simple interface focusing on core operations

## Supporting Components

### Input Validation (`internal/validator`)

**Purpose**: Validate and sanitize all user inputs

**Features**:

- Email format validation with RFC compliance
- Password strength requirements (length, complexity)
- Email normalization (case-insensitive, whitespace trimming)
- Clear, actionable error messages

### Security Utilities (`pkg/hash`, `pkg/token`)

**Purpose**: Cryptographic operations and token management

**Features**:

- bcrypt password hashing with secure defaults
- JWT token generation and validation
- Configurable token expiration
- Interface-based design for algorithm flexibility

### Middleware (`internal/middleware`)

**Purpose**: Cross-cutting concerns across HTTP requests

**Features**:

- Structured request logging with timing
- Authentication enforcement for protected routes
- Consistent JSON response formatting

## Data Flow

### User Registration Flow

```markdown
HTTP Request → Input Validation → Password Hashing → Store User → Generate JWT → Response
```

### User Authentication Flow  

```markdown
HTTP Request → Input Validation → Retrieve User → Password Verification → Generate JWT → Response
```

### Protected Route Access

```markdown
HTTP Request → Extract JWT → Validate Token → Execute Handler → Response
```

## Security Architecture

### Defense in Depth

1. **Input Layer**: Comprehensive validation, sanitization
2. **Authentication Layer**: JWT token validation, secure defaults
3. **Storage Layer**: bcrypt hashing, no plaintext passwords
4. **Transport Layer**: HTTPS recommended for production

### Security Controls

- **Password Security**: bcrypt with secure cost factor
- **Token Security**: HMAC-SHA256 signed JWTs with expiration
- **Input Security**: Validation against injection attacks
- **Error Security**: No sensitive information in error messages

## Scalability Design

### Stateless Architecture

- No server-side session storage
- JWT tokens carry all necessary user context
- Horizontal scaling friendly

### Performance Optimizations

- In-memory user store for development (easily swappable)
- Efficient bcrypt cost factor balancing security and performance
- Minimal memory allocations in hot paths

### Monitoring & Observability

- Health check endpoints for load balancer integration
- Structured logging with request correlation
- Performance metrics (request duration, status codes)

## Extension Points

### Database Integration

Replace `internal/store/memory` with database implementation:

```go
type PostgresUserStore struct {
    db *sql.DB
}

func (p *PostgresUserStore) Create(ctx context.Context, user *model.User) error {
    // Database implementation
}
```

### Additional Authentication Methods

Extend `internal/service/auth.go`:

```go
func (a *AuthService) SigninWithOAuth(provider string, token string) (string, error) {
    // OAuth implementation
}
```

### Enhanced Security

Add to `internal/middleware`:

```go
func RateLimitMiddleware(requests int, window time.Duration) func(http.Handler) http.Handler {
    // Rate limiting implementation
}
```

## Configuration Management

### Environment-Based Configuration

- `JWT_SECRET`: Token signing key
- `TOKEN_TTL_SECONDS`: Token expiration time  
- `HTTP_ADDR`: Server listen address

### Configuration Loading

- Secure defaults for development
- Environment variable override
- Validation of required settings

## Testing Strategy

### Unit Testing

- Each layer tested in isolation
- Mock dependencies for focused testing
- High coverage on critical security paths

### Integration Testing

- Full HTTP request/response cycle testing
- Real authentication flows
- Error scenario validation

### Security Testing

- Input validation edge cases
- Token security validation
- Authentication bypass attempts

This architecture provides a solid foundation for a production identity service while maintaining simplicity and clarity in the codebase.
