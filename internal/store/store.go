package store

import (
	"context"

	"github.com/coinbase/identity-service/internal/model"
)

type UserStore interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}
