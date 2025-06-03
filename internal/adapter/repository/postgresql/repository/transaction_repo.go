package repository

import (
	"errors"

	"github.com/hydr0g3nz/wallet_topup_system/internal/adapter/repository/postgresql/model"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/entity"
	errs "github.com/hydr0g3nz/wallet_topup_system/internal/domain/error"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}
func getQueryFromTrancsactionFilter(tx *gorm.DB, filter *entity.TransactionFilter) *gorm.DB {
	if filter == nil {
		return tx
	}
	if filter.ID != nil {
		tx = tx.Where("id = ?", filter.ID)
	}
	if filter.PaymentMethod != nil {
		tx = tx.Where("payment_method = ?", filter.PaymentMethod.String())
	}
	if filter.Status != nil {
		tx = tx.Where("status = ?", filter.Status.String())
	}
	if filter.Amount != nil {
		tx = tx.Where("amount = ?", filter.Amount.Amount())
	}
	if filter.ExpiredAt != nil {
		tx = tx.Where("expires_at <= ?", *filter.ExpiredAt)
	}
	return tx
}

func (r *TransactionRepository) FindAll(filter *entity.TransactionFilter) ([]entity.Transaction, error) {
	var transactionModels []model.Transaction
	query := r.db.Model(&model.Transaction{})
	query = getQueryFromTrancsactionFilter(query, filter)
	if err := query.Find(&transactionModels).Error; err != nil {
		return nil, err
	}
	transactions := make([]entity.Transaction, len(transactionModels))
	for i, tm := range transactionModels {
		t, err := tm.ToDomain()
		if err != nil {
			return nil, err
		}
		transactions[i] = *t
	}

	return transactions, nil
}

func (r *TransactionRepository) FindById(id uint) (*entity.Transaction, error) {
	var transactionModel model.Transaction
	if err := r.db.First(&transactionModel, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}
	t, err := transactionModel.ToDomain()
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *TransactionRepository) Create(transaction entity.Transaction) (uint, error) {
	transactionModel := model.CreateTransactionFromDomain(transaction)
	if err := r.db.Create(&transactionModel).Error; err != nil {
		return 0, err
	}
	return transactionModel.ID, nil
}
func (r *TransactionRepository) Update(filter *entity.TransactionFilter, transaction entity.Transaction) error {
	query := r.db.Model(&model.Transaction{})
	query = getQueryFromTrancsactionFilter(query, filter)
	if result := query.Updates(transaction.ToNotEmptyValueMap()); result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return errs.ErrNotFound
	}
	return nil
}
