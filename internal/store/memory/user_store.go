package memory

import (
	"context"
	"sync"
	"time"

	"github.com/coinbase/identity-service/internal/model"
	"github.com/google/uuid"
)

type UserStore struct {
	mu    sync.RWMutex
	users map[string]*model.User
}

func NewUserStore() *UserStore {
	return &UserStore{users: make(map[string]*model.User)}
}

func (s *UserStore) Create(_ context.Context, u *model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	u.ID = uuid.New()
	u.CreatedAt = now
	u.UpdatedAt = now
	s.users[u.Email] = u
	return nil
}

func (s *UserStore) GetByEmail(_ context.Context, email string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if u, ok := s.users[email]; ok {
		return u, nil
	}
	return nil, nil
}
