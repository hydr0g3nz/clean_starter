package repository

import (
	"errors"

	"github.com/hydr0g3nz/wallet_topup_system/internal/adapter/repository/postgresql/model"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/entity"
	errs "github.com/hydr0g3nz/wallet_topup_system/internal/domain/error"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}
func getQueryFromUserFilter(tx *gorm.DB, filter *entity.UserFilter) *gorm.DB {
	if filter == nil {
		return tx
	}
	if filter.FirstName != nil {
		tx = tx.Where("first_name = ?", *filter.FirstName)
	}
	if filter.LastName != nil {
		tx = tx.Where("last_name = ?", *filter.LastName)
	}
	if filter.Email != nil {
		tx = tx.Where("email = ?", *filter.Email)
	}
	if filter.Phone != nil {
		tx = tx.Where("phone = ?", *filter.Phone)
	}
	return tx
}

func (r *UserRepository) FindAll(userFilter *entity.UserFilter) ([]entity.User, error) {
	var userModels []model.User
	query := r.db.Model(&model.User{})
	query = getQueryFromUserFilter(query, userFilter)
	if err := query.Find(&userModels).Error; err != nil {
		return nil, err
	}
	users := make([]entity.User, len(userModels))
	for i, um := range userModels {
		users[i] = um.ToDomain()
	}
	return users, nil
}

func (r *UserRepository) FindById(id uint) (entity.User, error) {
	var userModel model.User
	if err := r.db.First(&userModel, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, errs.ErrNotFound
		}
		return entity.User{}, err
	}

	return userModel.ToDomain(), nil
}

func (r *UserRepository) Create(user entity.User) error {
	userModel := model.CreateUserFromDomain(user)
	return r.db.Create(&userModel).Error
}
func (r *UserRepository) Update(user entity.User) error {
	return r.db.Model(&model.User{}).Where("id = ?", user.ID).Updates(user.ToNotEmptyValueMap()).Error
}
