package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr  string
	JWTSecret string
	TokenTTL  time.Duration
}

func Load() Config {
	ttl, err := strconv.Atoi(getEnv("TOKEN_TTL_SECONDS", "900"))
	if err != nil {
		log.Fatalf("invalid TOKEN_TTL_SECONDS: %v", err)
	}
	return Config{
		HTTPAddr:  getEnv("HTTP_ADDR", ":8080"),
		JWTSecret: getEnvOrPanic("JWT_SECRET"),
		TokenTTL:  time.Duration(ttl) * time.Second,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
func getEnvOrPanic(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing required env %s", key)
	}
	return v
}
