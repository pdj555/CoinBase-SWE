package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Password  string // bcrypt hash
	CreatedAt time.Time
	UpdatedAt time.Time
}
