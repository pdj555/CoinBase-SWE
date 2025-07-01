package memory

import (
	"context"
	"testing"

	"github.com/coinbase/identity-service/internal/model"
)

func TestUserStore_Create(t *testing.T) {
	store := NewUserStore()
	ctx := context.Background()

	user := &model.User{
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	err := store.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Check that ID and timestamps were set
	if user.ID.String() == "" {
		t.Error("Create() should set user ID")
	}

	if user.CreatedAt.IsZero() {
		t.Error("Create() should set CreatedAt")
	}

	if user.UpdatedAt.IsZero() {
		t.Error("Create() should set UpdatedAt")
	}
}

func TestUserStore_GetByEmail(t *testing.T) {
	store := NewUserStore()
	ctx := context.Background()

	// Create a user
	user := &model.User{
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	err := store.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Retrieve the user
	retrieved, err := store.GetByEmail(ctx, user.Email)
	if err != nil {
		t.Fatalf("GetByEmail() failed: %v", err)
	}

	if retrieved == nil {
		t.Fatal("GetByEmail() returned nil user")
	}

	if retrieved.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, retrieved.Email)
	}

	if retrieved.Password != user.Password {
		t.Errorf("Expected password %s, got %s", user.Password, retrieved.Password)
	}
}

func TestUserStore_GetByEmailNotFound(t *testing.T) {
	store := NewUserStore()
	ctx := context.Background()

	// Try to retrieve non-existent user
	user, err := store.GetByEmail(ctx, "nonexistent@example.com")
	if err != nil {
		t.Fatalf("GetByEmail() failed: %v", err)
	}

	if user != nil {
		t.Error("GetByEmail() should return nil for non-existent user")
	}
}

func TestUserStore_CreateDuplicate(t *testing.T) {
	store := NewUserStore()
	ctx := context.Background()

	// Create first user
	user1 := &model.User{
		Email:    "test@example.com",
		Password: "hashedpassword1",
	}

	err := store.Create(ctx, user1)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Create second user with same email
	user2 := &model.User{
		Email:    "test@example.com",
		Password: "hashedpassword2",
	}

	err = store.Create(ctx, user2)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// The second user should overwrite the first
	retrieved, err := store.GetByEmail(ctx, "test@example.com")
	if err != nil {
		t.Fatalf("GetByEmail() failed: %v", err)
	}

	if retrieved.Password != "hashedpassword2" {
		t.Error("Second user should overwrite the first")
	}
}
