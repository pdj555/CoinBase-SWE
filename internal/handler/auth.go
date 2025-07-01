package handler

import (
	"encoding/json"
	"net/http"

	"github.com/coinbase/identity-service/internal/service"
	"github.com/coinbase/identity-service/internal/validator"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(a *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: a}
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req validator.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	token, err := h.auth.Signup(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandler) Signin(w http.ResponseWriter, r *http.Request) {
	var req validator.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	token, err := h.auth.Signin(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusUnauthorized)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandler) Me(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
