package repository

import (
	"fmt"

	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/repository"
	"gorm.io/gorm"
)

type DBTransactionRepository struct {
	db *gorm.DB
}

func NewDBTransactionRepository(db *gorm.DB) *DBTransactionRepository {
	return &DBTransactionRepository{db: db}
}
func (d *DBTransactionRepository) DoInTransaction(fn func(repo repository.Repository) error) error {
	// เริ่ม transaction
	tx := d.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	fmt.Println("start transaction")

	repoCtx := &RepositoryTransaction{
		transactionRepo: NewTransactionRepository(tx),
		walletRepo:      NewWalletRepository(tx),
	}
	// ทำงานภายใน transaction
	if err := fn(repoCtx); err != nil {
		tx.Rollback()
		return err
	}

	// commit transaction
	return tx.Commit().Error
}

type RepositoryTransaction struct {
	transactionRepo repository.TransactionRepository
	walletRepo      repository.WalletRepository
	userRepo        repository.UserRepository
}

func NewRepositoryTransaction(
	transactionRepo repository.TransactionRepository,
	walletRepo repository.WalletRepository,
	userRepo repository.UserRepository,
) *RepositoryTransaction {
	return &RepositoryTransaction{
		transactionRepo: transactionRepo,
		walletRepo:      walletRepo,
		userRepo:        userRepo,
	}
}
func (r *RepositoryTransaction) UserRepository() repository.UserRepository {
	return r.userRepo
}
func (r *RepositoryTransaction) WalletRepository() repository.WalletRepository {
	return r.walletRepo
}
func (r *RepositoryTransaction) TransactionRepository() repository.TransactionRepository {
	return r.transactionRepo
}
