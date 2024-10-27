package db

import (
	"github.com/cloud9cloud9/go-grpc-todo/auth-svc/internal/config"
	"github.com/cloud9cloud9/go-grpc-todo/auth-svc/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Repository struct {
	DB *gorm.DB
}

func ConnectToPostgreSQL(cfg *config.Config) Repository {
	dsn := cfg.UrlDB

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err)
	}

	err = db.AutoMigrate(
		&domain.User{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database. Error: %v", err)
	}

	return Repository{
		DB: db,
	}
}
