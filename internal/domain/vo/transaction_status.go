package vo

import (
	"strings"

	errs "github.com/hydr0g3nz/wallet_topup_system/internal/domain/error"
)

type TransactionStatus string

const (
	StatusVerified  TransactionStatus = "verified"
	StatusCompleted TransactionStatus = "completed"
	StatusFailed    TransactionStatus = "failed"
	StatusExpired   TransactionStatus = "expired"
)

func (s TransactionStatus) Valid() bool {
	switch s {
	case StatusVerified, StatusCompleted, StatusFailed, StatusExpired:
		return true
	default:
		return false
	}
}

func NewTransactionStatus(status string) (TransactionStatus, error) {
	s := TransactionStatus(strings.ToLower(status))
	if !s.Valid() {
		return "", errs.ErrInvalidTransactionStatus
	}
	return s, nil
}

func (s TransactionStatus) String() string {
	return string(s)
}
