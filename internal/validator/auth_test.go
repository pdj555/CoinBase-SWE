package validator

import (
	"testing"
)

func TestAuthRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     AuthRequest
		wantErr error
	}{
		{
			name: "valid request",
			req: AuthRequest{
				Email:    "test@coinbase.com",
				Password: "password123",
			},
			wantErr: nil,
		},
		{
			name: "email normalization",
			req: AuthRequest{
				Email:    "  Test@CoinBase.COM  ",
				Password: "password123",
			},
			wantErr: nil,
		},
		{
			name: "empty email",
			req: AuthRequest{
				Email:    "",
				Password: "password123",
			},
			wantErr: ErrEmailRequired,
		},
		{
			name: "invalid email format",
			req: AuthRequest{
				Email:    "invalid-email",
				Password: "password123",
			},
			wantErr: ErrEmailInvalid,
		},
		{
			name: "empty password",
			req: AuthRequest{
				Email:    "test@coinbase.com",
				Password: "",
			},
			wantErr: ErrPasswordRequired,
		},
		{
			name: "password too short",
			req: AuthRequest{
				Email:    "test@coinbase.com",
				Password: "short",
			},
			wantErr: ErrPasswordTooShort,
		},
		{
			name: "password without numbers",
			req: AuthRequest{
				Email:    "test@coinbase.com",
				Password: "onlyletters",
			},
			wantErr: ErrPasswordTooWeak,
		},
		{
			name: "password without letters",
			req: AuthRequest{
				Email:    "test@coinbase.com",
				Password: "12345678",
			},
			wantErr: ErrPasswordTooWeak,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalEmail := tt.req.Email
			err := tt.req.Validate()

			if err != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check email normalization for valid cases
			if tt.wantErr == nil && tt.name == "email normalization" {
				expectedEmail := "test@coinbase.com"
				if tt.req.Email != expectedEmail {
					t.Errorf("Email not normalized: got %s, want %s", tt.req.Email, expectedEmail)
				}
			}

			_ = originalEmail // Avoid unused variable warning
		})
	}
}
