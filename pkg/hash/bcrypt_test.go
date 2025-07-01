package hash

import (
	"testing"
)

func TestBcrypt_Hash(t *testing.T) {
	b := Bcrypt{}
	password := "testpassword123"

	hash, err := b.Hash(password)
	if err != nil {
		t.Fatalf("Hash() failed: %v", err)
	}

	if hash == "" {
		t.Error("Hash() returned empty string")
	}

	if hash == password {
		t.Error("Hash() returned the original password")
	}
}

func TestBcrypt_Compare(t *testing.T) {
	b := Bcrypt{}
	password := "testpassword123"
	wrongPassword := "wrongpassword"

	hash, err := b.Hash(password)
	if err != nil {
		t.Fatalf("Hash() failed: %v", err)
	}

	// Test correct password
	if !b.Compare(hash, password) {
		t.Error("Compare() failed for correct password")
	}

	// Test wrong password
	if b.Compare(hash, wrongPassword) {
		t.Error("Compare() should fail for wrong password")
	}

	// Test empty password
	if b.Compare(hash, "") {
		t.Error("Compare() should fail for empty password")
	}
}
