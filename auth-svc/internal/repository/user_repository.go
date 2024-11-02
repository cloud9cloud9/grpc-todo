package repository

import (
	"github.com/cloud9cloud9/go-grpc-todo/auth-svc/internal/domain"
	"gorm.io/gorm"
)

type AuthPostgres struct {
	db *gorm.DB
}

func NewAuthPostgres(db *gorm.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (ap *AuthPostgres) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	result := ap.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (ap *AuthPostgres) CreateUser(user *domain.User) error {
	if err := ap.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (ap *AuthPostgres) FindByID(id int64) (*domain.User, error) {
	var user domain.User
	result := ap.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
