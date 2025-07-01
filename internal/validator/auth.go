package validator

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrEmailRequired    = errors.New("email is required")
	ErrEmailInvalid     = errors.New("email format is invalid")
	ErrPasswordRequired = errors.New("password is required")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrPasswordTooWeak  = errors.New("password must contain letters and numbers")
)

// emailRegex is a basic email validation regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var hasLetter = regexp.MustCompile(`[a-zA-Z]`)
var hasNumber = regexp.MustCompile(`[0-9]`)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *AuthRequest) Validate() error {
	// Email validation
	if strings.TrimSpace(a.Email) == "" {
		return ErrEmailRequired
	}

	a.Email = strings.TrimSpace(strings.ToLower(a.Email))

	if !emailRegex.MatchString(a.Email) {
		return ErrEmailInvalid
	}

	// Password validation
	if a.Password == "" {
		return ErrPasswordRequired
	}

	if len(a.Password) < 8 {
		return ErrPasswordTooShort
	}

	if !hasLetter.MatchString(a.Password) || !hasNumber.MatchString(a.Password) {
		return ErrPasswordTooWeak
	}

	return nil
}
