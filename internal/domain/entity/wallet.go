package entity

import "github.com/hydr0g3nz/wallet_topup_system/internal/domain/vo"

// Wallet represents the wallets table (1-to-1 with User)
type Wallet struct {
	ID      uint
	Balance vo.Money
}

func (w Wallet) ToNotEmptyValueMap() map[string]interface{} {
	result := make(map[string]interface{})
	if !w.Balance.IsZero() {
		result["balance"] = w.Balance.Amount()
	}
	return result
}

type WalletFilter struct {
	Balance *vo.Money
}
