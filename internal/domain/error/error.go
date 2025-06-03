package errs

import "errors"

var ErrNegativeAmount = errors.New("amount cannot be negative")
var ErrInsufficientBalance = errors.New("insufficient balance")
var ErrExpiredTransaction = errors.New("transaction expired")
var ErrNotFound = errors.New("not found")
var ErrTransactionNotVerified = errors.New("transaction is not in 'verified' status")
var ErrAmountExceedsLimit = errors.New("amount exceeds maximum limit")
var ErrInvalidPaymentMethod = errors.New("invalid payment method")
var ErrInvalidTransactionStatus = errors.New("invalid transaction status")
