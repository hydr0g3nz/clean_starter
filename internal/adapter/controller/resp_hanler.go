package controller

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	errs "github.com/hydr0g3nz/wallet_topup_system/internal/domain/error"
)

type successResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// SuccessResp builds a success response
func SuccessResp(c *fiber.Ctx, status int, message string, data any) error {
	return c.Status(status).JSON(successResponse{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

// HandleError builds an appropriate Fiber error response based on the domain error
func HandleError(c *fiber.Ctx, err error) error {
	var statusCode int
	var message string

	switch {
	case errors.Is(err, errs.ErrNegativeAmount):
		statusCode = http.StatusBadRequest
		message = "Amount cannot be negative"
	case errors.Is(err, errs.ErrInsufficientBalance):
		statusCode = http.StatusBadRequest
		message = "Insufficient balance"
	case errors.Is(err, errs.ErrExpiredTransaction):
		statusCode = http.StatusBadRequest
		message = "Transaction expired"
	case errors.Is(err, errs.ErrNotFound):
		statusCode = http.StatusNotFound
		message = "Not found"
	case errors.Is(err, errs.ErrTransactionNotVerified):
		statusCode = http.StatusBadRequest
		message = "Transaction is not in 'verified' status"
	case errors.Is(err, errs.ErrAmountExceedsLimit):
		statusCode = http.StatusBadRequest
		message = "Amount exceeds maximum limit"
	case errors.Is(err, errs.ErrInvalidPaymentMethod):
		statusCode = http.StatusBadRequest
		message = "Invalid payment method"
	case errors.Is(err, errs.ErrInvalidTransactionStatus):
		statusCode = http.StatusBadRequest
		message = "Invalid transaction status"
	default:
		statusCode = http.StatusInternalServerError
		message = "Something went wrong"
	}

	return c.Status(statusCode).JSON(ErrorResponse{
		Status:  statusCode,
		Message: message,
	})
}
