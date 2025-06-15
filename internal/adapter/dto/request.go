package dto

import "time"

// VerifyRequest represents the input data for verifying a top-up request
type VerifyRequest struct {
	UserID        uint    `json:"user_id"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
}

// VerifyResponse represents the output data for verifying a top-up request
type VerifyResponse struct {
	TransactionID uint      `json:"transaction_id"`
	UserID        uint      `json:"user_id"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	ExpiresAt     time.Time `json:"expires_at"`
}

// ConfirmRequest represents the input data for confirming a top-up transaction
type ConfirmRequest struct {
	TransactionID uint `json:"transaction_id"`
}

// ConfirmResponse represents the output data for confirming a top-up transaction
type ConfirmResponse struct {
	TransactionID uint    `json:"transaction_id"`
	UserID        uint    `json:"user_id"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
	Balance       float64 `json:"balance"`
}

// RegisterRequest represents user registration request DTO
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=candidate company_hr admin"`
}

// LoginRequest represents user login request DTO
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateProfileRequest represents user profile update request DTO
type UpdateProfileRequest struct {
	Email string `json:"email,omitempty" validate:"omitempty,email"`
}

// ChangePasswordRequest represents password change request DTO
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// UserResponse represents user data in responses
type UserResponse struct {
	ID            int        `json:"id"`
	Email         string     `json:"email"`
	Role          string     `json:"role"`
	IsActive      bool       `json:"is_active"`
	EmailVerified bool       `json:"email_verified"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
}

// LoginResponse represents login response with token
type LoginResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

// UsersListResponse represents paginated users list response
type UsersListResponse struct {
	Users  []*UserResponse `json:"users"`
	Total  int             `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}
