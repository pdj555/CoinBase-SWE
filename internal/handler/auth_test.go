package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coinbase/identity-service/internal/service"
	"github.com/coinbase/identity-service/internal/store/memory"
	"github.com/coinbase/identity-service/pkg/hash"
	"github.com/coinbase/identity-service/pkg/token"
)

func setupAuthHandler() *AuthHandler {
	userStore := memory.NewUserStore()
	hasher := hash.Bcrypt{}
	tokens := token.NewJWTManager("test-secret-key", 15*time.Minute)
	authSvc := service.NewAuthService(userStore, hasher, tokens)
	return NewAuthHandler(authSvc)
}

func TestAuthHandler_Signup(t *testing.T) {
	handler := setupAuthHandler()

	requestBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Signup(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["token"] == "" {
		t.Error("Response should contain a token")
	}
}

func TestAuthHandler_SignupInvalidJSON(t *testing.T) {
	handler := setupAuthHandler()

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Signup(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_SignupDuplicateEmail(t *testing.T) {
	handler := setupAuthHandler()

	requestBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}

	body, _ := json.Marshal(requestBody)

	// First signup
	req1 := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	handler.Signup(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("First signup should succeed, got status %d", w1.Code)
	}

	// Second signup with same email
	req2 := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	handler.Signup(w2, req2)

	if w2.Code != http.StatusBadRequest {
		t.Errorf("Second signup should fail with 400, got status %d", w2.Code)
	}
}

func TestAuthHandler_Signin(t *testing.T) {
	handler := setupAuthHandler()

	requestBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}

	body, _ := json.Marshal(requestBody)

	// First signup
	req1 := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	handler.Signup(w1, req1)

	// Then signin
	req2 := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	handler.Signin(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w2.Code)
	}

	var response map[string]string
	err := json.NewDecoder(w2.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["token"] == "" {
		t.Error("Response should contain a token")
	}
}

func TestAuthHandler_SigninInvalidCredentials(t *testing.T) {
	handler := setupAuthHandler()

	// Signup first
	signupBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(signupBody)
	req1 := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	handler.Signup(w1, req1)

	// Try signin with wrong password (but valid format)
	signinBody := map[string]string{
		"email":    "test@example.com",
		"password": "wrongpass123",
	}
	body, _ = json.Marshal(signinBody)
	req2 := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	handler.Signin(w2, req2)

	if w2.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w2.Code)
	}
}

func TestAuthHandler_SigninInvalidPassword(t *testing.T) {
	handler := setupAuthHandler()

	// Try signin with password that doesn't meet validation rules
	signinBody := map[string]string{
		"email":    "test@example.com",
		"password": "weak", // Too short and no numbers
	}
	body, _ := json.Marshal(signinBody)
	req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Signin(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_Me(t *testing.T) {
	handler := setupAuthHandler()

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	w := httptest.NewRecorder()
	handler.Me(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "ok" {
		t.Error("Response should contain status: ok")
	}
}
