package server

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/coinbase/identity-service/internal/handler"
	"github.com/coinbase/identity-service/internal/middleware"
	"github.com/coinbase/identity-service/internal/service"
	"github.com/coinbase/identity-service/pkg/token"
)

func NewRouter(authSvc *service.AuthService, tm token.Manager) *mux.Router {
	authHandler := handler.NewAuthHandler(authSvc)
	healthHandler := handler.NewHealthHandler()

	r := mux.NewRouter()

	// Global middleware
	r.Use(middleware.LoggingMiddleware)
	r.Use(jsonMiddleware)

	// Health check endpoints
	r.HandleFunc("/health", healthHandler.Health).Methods(http.MethodGet)
	r.HandleFunc("/ready", healthHandler.Ready).Methods(http.MethodGet)

	// Authentication endpoints
	r.HandleFunc("/signup", authHandler.Signup).Methods(http.MethodPost)
	r.HandleFunc("/signin", authHandler.Signin).Methods(http.MethodPost)

	// Protected endpoints
	r.Handle("/me", authMiddleware(tm, authHandler.Me)).Methods(http.MethodGet)

	return r
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(tm token.Manager, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		fields := strings.Fields(auth)
		if len(fields) != 2 || fields[0] != "Bearer" {
			http.Error(w, `{"error":"missing token"}`, http.StatusUnauthorized)
			return
		}
		claims, err := tm.Verify(fields[1])
		if err != nil {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}
		// attach claims to context if needed
		r = r.WithContext(r.Context())
		_ = claims
		next.ServeHTTP(w, r)
	}
}
