package token

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTManager_Generate(t *testing.T) {
	secret := "test-secret-key"
	ttl := 15 * time.Minute
	jm := NewJWTManager(secret, ttl)

	userID := uuid.New()
	email := "test@example.com"

	token, err := jm.Generate(userID, email)
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	if token == "" {
		t.Error("Generate() returned empty token")
	}
}

func TestJWTManager_Verify(t *testing.T) {
	secret := "test-secret-key"
	ttl := 15 * time.Minute
	jm := NewJWTManager(secret, ttl)

	userID := uuid.New()
	email := "test@example.com"

	// Generate a token
	token, err := jm.Generate(userID, email)
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	// Verify the token
	claims, err := jm.Verify(token)
	if err != nil {
		t.Fatalf("Verify() failed: %v", err)
	}

	if claims.UserID != userID.String() {
		t.Errorf("Expected UserID %s, got %s", userID.String(), claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("Expected Email %s, got %s", email, claims.Email)
	}
}

func TestJWTManager_VerifyInvalidToken(t *testing.T) {
	secret := "test-secret-key"
	ttl := 15 * time.Minute
	jm := NewJWTManager(secret, ttl)

	// Test invalid token
	_, err := jm.Verify("invalid-token")
	if err == nil {
		t.Error("Verify() should fail for invalid token")
	}

	// Test empty token
	_, err = jm.Verify("")
	if err == nil {
		t.Error("Verify() should fail for empty token")
	}
}

func TestJWTManager_VerifyExpiredToken(t *testing.T) {
	secret := "test-secret-key"
	ttl := 1 * time.Millisecond // Very short TTL
	jm := NewJWTManager(secret, ttl)

	userID := uuid.New()
	email := "test@example.com"

	// Generate a token
	token, err := jm.Generate(userID, email)
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	// Verify should fail for expired token
	_, err = jm.Verify(token)
	if err == nil {
		t.Error("Verify() should fail for expired token")
	}
}
