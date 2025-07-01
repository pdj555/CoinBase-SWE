package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/coinbase/identity-service/internal/store/memory"
	"github.com/coinbase/identity-service/pkg/hash"
	"github.com/coinbase/identity-service/pkg/token"
)

func setupAuthService() *AuthService {
	userStore := memory.NewUserStore()
	hasher := hash.Bcrypt{}
	tokens := token.NewJWTManager("test-secret-key", 15*time.Minute)
	return NewAuthService(userStore, hasher, tokens)
}

func TestAuthService_Signup(t *testing.T) {
	auth := setupAuthService()
	ctx := context.Background()

	email := "test@example.com"
	password := "password123"

	token, err := auth.Signup(ctx, email, password)
	if err != nil {
		t.Fatalf("Signup() failed: %v", err)
	}

	if token == "" {
		t.Error("Signup() should return a token")
	}

	// Verify token is valid
	if len(strings.Split(token, ".")) != 3 {
		t.Error("Token should be a valid JWT with 3 parts")
	}
}

func TestAuthService_SignupDuplicateEmail(t *testing.T) {
	auth := setupAuthService()
	ctx := context.Background()

	email := "test@example.com"
	password := "password123"

	// First signup should succeed
	_, err := auth.Signup(ctx, email, password)
	if err != nil {
		t.Fatalf("First Signup() failed: %v", err)
	}

	// Second signup with same email should fail
	_, err = auth.Signup(ctx, email, password)
	if err != ErrUserExists {
		t.Errorf("Expected ErrUserExists, got %v", err)
	}
}

func TestAuthService_Signin(t *testing.T) {
	auth := setupAuthService()
	ctx := context.Background()

	email := "test@example.com"
	password := "password123"

	// First signup
	_, err := auth.Signup(ctx, email, password)
	if err != nil {
		t.Fatalf("Signup() failed: %v", err)
	}

	// Then signin
	token, err := auth.Signin(ctx, email, password)
	if err != nil {
		t.Fatalf("Signin() failed: %v", err)
	}

	if token == "" {
		t.Error("Signin() should return a token")
	}
}

func TestAuthService_SigninInvalidEmail(t *testing.T) {
	auth := setupAuthService()
	ctx := context.Background()

	// Try to signin with non-existent email
	_, err := auth.Signin(ctx, "nonexistent@example.com", "password123")
	if err != ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestAuthService_SigninInvalidPassword(t *testing.T) {
	auth := setupAuthService()
	ctx := context.Background()

	email := "test@example.com"
	password := "password123"
	wrongPassword := "wrongpassword"

	// First signup
	_, err := auth.Signup(ctx, email, password)
	if err != nil {
		t.Fatalf("Signup() failed: %v", err)
	}

	// Try signin with wrong password
	_, err = auth.Signin(ctx, email, wrongPassword)
	if err != ErrInvalidCreds {
		t.Errorf("Expected ErrInvalidCreds, got %v", err)
	}
}

func TestAuthService_SigninEmptyPassword(t *testing.T) {
	auth := setupAuthService()
	ctx := context.Background()

	email := "test@example.com"
	password := "password123"

	// First signup
	_, err := auth.Signup(ctx, email, password)
	if err != nil {
		t.Fatalf("Signup() failed: %v", err)
	}

	// Try signin with empty password
	_, err = auth.Signin(ctx, email, "")
	if err != ErrInvalidCreds {
		t.Errorf("Expected ErrInvalidCreds, got %v", err)
	}
}
