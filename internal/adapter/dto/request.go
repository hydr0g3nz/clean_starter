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
