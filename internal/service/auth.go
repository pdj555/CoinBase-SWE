package service

import (
	"context"
	"errors"

	"github.com/coinbase/identity-service/internal/model"
	"github.com/coinbase/identity-service/internal/store"
	"github.com/coinbase/identity-service/pkg/hash"
	"github.com/coinbase/identity-service/pkg/token"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrInvalidCreds = errors.New("invalid credentials")
	ErrUserNotFound = errors.New("user not found")
)

type AuthService struct {
	users  store.UserStore
	hasher hash.Bcrypt
	tokens token.Manager
}

func NewAuthService(us store.UserStore, h hash.Bcrypt, t token.Manager) *AuthService {
	return &AuthService{users: us, hasher: h, tokens: t}
}

func (a *AuthService) Signup(ctx context.Context, email, password string) (string, error) {
	if existing, _ := a.users.GetByEmail(ctx, email); existing != nil {
		return "", ErrUserExists
	}
	hashPw, err := a.hasher.Hash(password)
	if err != nil {
		return "", err
	}
	u := &model.User{Email: email, Password: hashPw}
	if err := a.users.Create(ctx, u); err != nil {
		return "", err
	}
	return a.tokens.Generate(u.ID, u.Email)
}

func (a *AuthService) Signin(ctx context.Context, email, password string) (string, error) {
	u, err := a.users.GetByEmail(ctx, email)
	if err != nil || u == nil {
		return "", ErrUserNotFound
	}
	if !a.hasher.Compare(u.Password, password) {
		return "", ErrInvalidCreds
	}
	return a.tokens.Generate(u.ID, u.Email)
}
