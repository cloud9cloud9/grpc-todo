package repository

import (
	"github.com/cloud9cloud9/go-grpc-todo/auth-svc/internal/domain"
	"gorm.io/gorm"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Repository struct {
	UserRepository
}

type UserRepository interface {
	FindByEmail(email string) (*domain.User, error)
	CreateUser(user *domain.User) error
	FindByID(id int64) (*domain.User, error)
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		UserRepository: NewAuthPostgres(db),
	}
}
