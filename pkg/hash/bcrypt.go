package hash

import (
	"golang.org/x/crypto/bcrypt"
)

type Bcrypt struct{}

// Hash returns bcrypt hash of the password.
func (Bcrypt) Hash(pw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(bytes), err
}

// Compare verifies a bcrypt‑hashed password with its possible plaintext equivalent.
func (Bcrypt) Compare(hashed, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}
