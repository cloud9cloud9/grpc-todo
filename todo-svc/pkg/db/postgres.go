package db

import (
	"fmt"
	"github.com/cloud9cloud9/go-grpc-todo/todo-svc/internal/config"
	"github.com/cloud9cloud9/go-grpc-todo/todo-svc/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func ConnectToPostgreSQL(cfg *config.Config) (db *gorm.DB) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.PostgresUser, cfg.Database.PostgresDB,
		cfg.Database.PostgresPassword, cfg.Database.PostgresSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err)
	}

	err = db.AutoMigrate(
		&domain.TodoList{},
		&domain.UsersList{},
		&domain.TodoItem{},
		&domain.ListsItem{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database. Error: %v", err)
	}

	return db
}
