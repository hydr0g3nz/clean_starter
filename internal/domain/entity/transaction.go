package entity

import (
	"time"

	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/vo"
)

// Transaction represents the transactions table
type Transaction struct {
	ID            uint                 `json:"id"`
	UserID        uint                 `json:"user_id"`
	Amount        vo.Money             `json:"amount"`
	PaymentMethod vo.PaymentMethod     `json:"payment_method"`
	Status        vo.TransactionStatus `json:"status"`
	ExpiresAt     time.Time            `json:"expires_at"`
}

func NewTransaction(UserID uint, amount float64, paymentMethod string, status string, expiresAt time.Time) (Transaction, error) {

	newPaymentMethod, err := vo.NewPaymentMethod(paymentMethod)
	if err != nil {
		return Transaction{}, err
	}
	newStatus, err := vo.NewTransactionStatus(status)
	if err != nil {
		return Transaction{}, err
	}
	newAmount, err := vo.NewMoney(amount)
	if err != nil {
		return Transaction{}, err
	}
	return Transaction{
		UserID:        UserID,
		Amount:        newAmount,
		PaymentMethod: newPaymentMethod,
		Status:        newStatus,
		ExpiresAt:     expiresAt,
	}, nil
}
func (t Transaction) ToNotEmptyValueMap() map[string]interface{} {
	result := make(map[string]interface{})
	if !t.Amount.IsZero() {
		result["amount"] = t.Amount
	}
	if t.PaymentMethod != "" {
		result["payment_method"] = t.PaymentMethod.String()
	}
	if t.Status != "" {
		result["status"] = t.Status.String()
	}
	if !t.ExpiresAt.IsZero() {
		result["expires_at"] = t.ExpiresAt
	}
	return result
}

type TransactionFilter struct {
	ID            *uint
	PaymentMethod *vo.PaymentMethod
	Status        *vo.TransactionStatus
	Amount        *vo.Money
	ExpiredAt     *time.Time
}
