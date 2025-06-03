package model

import (
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/entity"
	"gorm.io/gorm"
)

// User represents the users table
type User struct {
	gorm.Model
	FirstName string `gorm:"size:50;not null"`
	LastName  string `gorm:"size:50;not null"`
	Email     string `gorm:"size:100;uniqueIndex;not null"`
	Password  string `gorm:"size:255;not null"`
	Phone     string `gorm:"size:20;not null"`
}

func CreateUserFromDomain(u entity.User) User {
	return User{
		Model:     gorm.Model{ID: u.ID},
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  u.Password,
		Phone:     u.Phone,
	}
}

func (u User) ToDomain() entity.User {
	return entity.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  u.Password,
		Phone:     u.Phone,
	}
}
