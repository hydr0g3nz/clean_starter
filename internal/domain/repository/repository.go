package repository

import "github.com/hydr0g3nz/wallet_topup_system/internal/domain/entity"

type Repository interface {
	UserRepository() UserRepository
	WalletRepository() WalletRepository
	TransactionRepository() TransactionRepository
}
type DBTransaction interface {
	DoInTransaction(fn func(repo Repository) error) error
}

type UserRepository interface {
	FindAll(*entity.UserFilter) ([]entity.User, error)
	FindById(id uint) (entity.User, error)
	Create(entity.User) error
}

type WalletRepository interface {
	Create(entity.Wallet) error
	Update(entity.Wallet) error
	FindById(id uint) (*entity.Wallet, error)
}

type TransactionRepository interface {
	FindAll(*entity.TransactionFilter) ([]entity.Transaction, error)
	FindById(id uint) (*entity.Transaction, error)
	Create(entity.Transaction) (uint, error)
	Update(*entity.TransactionFilter, entity.Transaction) error
}
