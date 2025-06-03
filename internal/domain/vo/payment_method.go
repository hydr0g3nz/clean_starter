package vo

import (
	"strings"

	errs "github.com/hydr0g3nz/wallet_topup_system/internal/domain/error"
)

type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
)

func (p PaymentMethod) Valid() bool {
	switch p {
	case PaymentMethodCreditCard:
		return true
	default:
		return false
	}
}

func NewPaymentMethod(method string) (PaymentMethod, error) {
	pm := PaymentMethod(strings.ToLower(method))
	if !pm.Valid() {
		return "", errs.ErrInvalidPaymentMethod
	}
	return pm, nil
}

func (p PaymentMethod) String() string {
	return string(p)
}
