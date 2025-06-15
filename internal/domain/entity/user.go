package entity

import (
	"time"

	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/vo"
)

// User represents the user domain entity
type User struct {
	ID            int         `json:"id"`
	Email         string      `json:"email"`
	PasswordHash  string      `json:"-"` // Hide password hash in JSON
	Role          vo.UserRole `json:"role"`
	IsActive      bool        `json:"is_active"`
	EmailVerified bool        `json:"email_verified"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	LastLoginAt   *time.Time  `json:"last_login_at,omitempty"`
}

func (u *User) IsValid() bool {
	if u.Email == "" || u.PasswordHash == "" || u.Role == "" {
		return false
	}
	return true
}
