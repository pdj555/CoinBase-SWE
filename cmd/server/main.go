package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"

	"github.com/coinbase/identity-service/internal/config"
	"github.com/coinbase/identity-service/internal/server"
	"github.com/coinbase/identity-service/internal/service"
	"github.com/coinbase/identity-service/internal/store/memory"
	"github.com/coinbase/identity-service/pkg/hash"
	"github.com/coinbase/identity-service/pkg/token"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	// ── infrastructure
	userStore := memory.NewUserStore()
	hasher := hash.Bcrypt{}
	tokens := token.NewJWTManager(cfg.JWTSecret, cfg.TokenTTL)

	// ── services
	authSvc := service.NewAuthService(userStore, hasher, tokens)

	// ── HTTP server
	r := server.NewRouter(authSvc, tokens)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("identity‑service listening on %s", cfg.HTTPAddr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-context.Background().Done()
}
