package repository

import (
	"errors"
	"fmt"
	"github.com/cloud9cloud9/go-grpc-todo/todo-svc/internal/domain"
	"gorm.io/gorm"
)

var (
	ErrTodoListNotFound  = errors.New("todo list not found")
	ErrUsersListNotFound = errors.New("users list not found")
)

type ListPostgres struct {
	db *gorm.DB
}

func NewTodoListPostgres(db *gorm.DB) *ListPostgres {
	return &ListPostgres{
		db: db,
	}
}

func (lp *ListPostgres) Create(userId int64, list *domain.TodoList) error {
	tx := lp.db.Begin()
	if err := tx.Create(list).Error; err != nil {
		tx.Rollback()
		return err
	}

	userList := domain.UsersList{
		UserId: userId,
		ListId: list.Id,
	}
	if err := tx.Create(&userList).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (lp *ListPostgres) GetAll(userId int64) ([]*domain.TodoList, error) {
	var userLists []domain.UsersList
	if err := lp.db.Where(&domain.UsersList{UserId: userId}).Find(&userLists).Error; err != nil {
		return nil, ErrUsersListNotFound
	}

	var todoLists []*domain.TodoList
	for _, userList := range userLists {
		var list *domain.TodoList
		if err := lp.db.Where(&domain.TodoList{Id: userList.ListId}).First(&list).Error; err == nil {
			todoLists = append(todoLists, list)
		}
	}
	return todoLists, nil
}

func (lp *ListPostgres) GetById(listId int64) (*domain.TodoList, error) {
	var list domain.TodoList
	if result := lp.db.Where(&domain.TodoList{Id: listId}).First(&list); result.Error != nil {
		return nil, result.Error
	}
	return &list, nil
}

func (lp *ListPostgres) Delete(listId int64) error {
	tx := lp.db.Begin()

	var list domain.TodoList
	if err := tx.Where(&domain.TodoList{Id: listId}).First(&list).Error; err != nil {
		tx.Rollback()
		return ErrTodoListNotFound
	}

	if err := tx.Where("list_id = ?", listId).Delete(&domain.TodoItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	var userList domain.UsersList
	if err := tx.Where(&domain.UsersList{ListId: listId}).First(&userList).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&userList).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&list).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (lp *ListPostgres) Update(listId int64, input *domain.TodoList) (*domain.TodoList, error) {
	var list domain.TodoList
	if err := lp.db.Where(&domain.TodoList{Id: listId}).First(&list).Error; err != nil {
		return nil, ErrTodoListNotFound
	}

	list.Title = input.Title
	if err := lp.db.Save(&list).Error; err != nil {
		return nil, err
	}

	return &list, nil
}

func (lp *ListPostgres) CheckUserAccessToList(userId int64, listId int64) error {
	var count int64

	if err := lp.db.Model(&domain.UsersList{}).
		Where("user_id = ? AND list_id = ?", userId, listId).
		Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("user does not have access to the list")
	}

	return nil
}
