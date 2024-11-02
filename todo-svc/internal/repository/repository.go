package repository

import (
	"github.com/cloud9cloud9/go-grpc-todo/todo-svc/internal/domain"
	"gorm.io/gorm"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type TodoList interface {
	Create(userId int64, list *domain.TodoList) error
	GetAll(userId int64) ([]*domain.TodoList, error)
	GetById(listId int64) (*domain.TodoList, error)
	Delete(listId int64) error
	Update(listId int64, input *domain.TodoList) (*domain.TodoList, error)
	CheckUserAccessToList(userId int64, listId int64) error
}

type TodoItem interface {
	Create(item *domain.TodoItem) error
	GetAll(listId int64) ([]*domain.TodoItem, error)
	GetById(itemId int64) (*domain.TodoItem, int64, error)
	Delete(itemId int64) error
	Update(input *domain.TodoItem) error
}

type Repository struct {
	TodoList
	TodoItem
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		TodoList: NewTodoListPostgres(db),
		TodoItem: NewTodoItemPostgres(db),
	}
}
