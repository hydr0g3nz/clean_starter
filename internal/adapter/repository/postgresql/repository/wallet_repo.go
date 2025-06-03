package repository

import (
	"github.com/hydr0g3nz/wallet_topup_system/internal/adapter/repository/postgresql/model"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/entity"
	"gorm.io/gorm"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Create(wallet entity.Wallet) error {
	walletModel := model.CreateWalletFromDomain(wallet)
	return r.db.Create(&walletModel).Error
}

func (r *WalletRepository) Update(wallet entity.Wallet) error {
	return r.db.Model(&model.Wallet{}).Where("id = ?", wallet.ID).Updates(wallet.ToNotEmptyValueMap()).Error
}
func (r *WalletRepository) FindById(id uint) (*entity.Wallet, error) {
	var walletModel model.Wallet
	if err := r.db.Take(&walletModel, id).Error; err != nil {
		return nil, err
	}
	w, err := walletModel.ToDomain()
	if err != nil {
		return nil, err
	}
	return &w, nil
}
