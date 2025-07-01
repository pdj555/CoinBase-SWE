# API Documentation

## Base URL

```markdown
http://localhost:8080
```

## Authentication

Protected endpoints require a JWT token in the Authorization header:

```markdown
Authorization: Bearer YOUR_JWT_TOKEN
```

## Public Endpoints

### Register User

Register a new user account.

**Endpoint**: `POST /signup`

**Request Body**:

```json
{
  "email": "user@example.com",
  "password": "securepass123"
}
```

**Validation Rules**:

- Email: Valid email format, automatically normalized to lowercase
- Password: Minimum 8 characters, must contain letters and numbers

**Success Response** (200):

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Responses**:

- `400` - Email already exists
- `400` - Invalid email format  
- `400` - Password too weak
- `400` - Invalid JSON

**Example**:

```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@coinbase.com",
    "password": "mypassword123"
  }'
```

---

### User Login

Authenticate an existing user.

**Endpoint**: `POST /signin`

**Request Body**:

```json
{
  "email": "user@example.com", 
  "password": "securepass123"
}
```

**Success Response** (200):

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Responses**:

- `401` - Invalid credentials
- `401` - User not found
- `400` - Invalid input format

**Example**:

```bash
curl -X POST http://localhost:8080/signin \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@coinbase.com",
    "password": "mypassword123"
  }'
```

## Protected Endpoints

### Get User Profile

Retrieve current user information.

**Endpoint**: `GET /me`

**Headers**:

```markdown
Authorization: Bearer YOUR_JWT_TOKEN
```

**Success Response** (200):

```json
{
  "status": "ok"
}
```

**Error Responses**:

- `401` - Missing token
- `401` - Invalid token
- `401` - Expired token

**Example**:

```bash
curl -X GET http://localhost:8080/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## Health Endpoints

### Service Health

Get service health status and metadata.

**Endpoint**: `GET /health`

**Success Response** (200):

```json
{
  "status": "healthy",
  "service": "identity-service", 
  "uptime": "2m30.123456789s",
  "timestamp": "2025-01-15T10:30:00Z",
  "version": "1.0.0"
}
```

**Example**:

```bash
curl http://localhost:8080/health
```

---

### Readiness Check

Check if service is ready to accept traffic.

**Endpoint**: `GET /ready`

**Success Response** (200):

```json
{
  "status": "ready"
}
```

**Example**:

```bash
curl http://localhost:8080/ready
```

## Error Handling

All error responses follow a consistent format:

```json
{
  "error": "descriptive error message"
}
```

### HTTP Status Codes

- `200` - Success
- `400` - Bad Request (validation errors, malformed JSON)
- `401` - Unauthorized (missing/invalid token, wrong credentials)
- `500` - Internal Server Error

### Common Error Messages

**Validation Errors**:

- `"email is required"`
- `"email format is invalid"`
- `"password is required"`
- `"password must be at least 8 characters"`
- `"password must contain letters and numbers"`

**Authentication Errors**:

- `"user already exists"`
- `"invalid credentials"`
- `"user not found"`
- `"missing token"`
- `"invalid token"`

## JWT Token Details

### Token Structure

The service uses HMAC-SHA256 signed JWT tokens containing:

```json
{
  "user_id": "uuid-string",
  "email": "user@example.com",
  "exp": 1642234567,
  "iat": 1642230967
}
```

### Token Expiration

- Default: 15 minutes (900 seconds)
- Configurable via `TOKEN_TTL_SECONDS` environment variable
- Tokens cannot be refreshed (obtain new token via `/signin`)

### Token Usage

Include in Authorization header with Bearer scheme:

```markdown
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Rate Limiting

Currently no rate limiting is implemented. For production use, consider:

- Adding rate limiting middleware
- Implementing per-IP request limits
- Adding authentication attempt limits

## CORS

No CORS headers are set by default. For browser applications, configure CORS middleware as needed.

## Content Type

All endpoints expect and return `application/json` content type.

## Request/Response Examples

### Complete Registration Flow

```bash
# 1. Register new user
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@example.com","password":"demopass123"}'

# Response:
# {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}

# 2. Use token to access protected resource
curl -X GET http://localhost:8080/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Response:
# {"status":"ok"}
```

### Complete Login Flow

```bash
# 1. Login existing user
curl -X POST http://localhost:8080/signin \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@example.com","password":"demopass123"}'

# Response:
# {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}

# 2. Access protected endpoint
curl -X GET http://localhost:8080/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Response:
# {"status":"ok"}
```
