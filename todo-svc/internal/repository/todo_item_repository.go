package repository

import (
	"errors"
	"github.com/cloud9cloud9/go-grpc-todo/todo-svc/internal/domain"
	"gorm.io/gorm"
)

var (
	ErrTodoItemNotFound = errors.New("todo item not found")
	ErrCreateTodoItem   = errors.New("failed to create todo item")
	ErrCreateListItem   = errors.New("failed to create list item")
	ErrListItemNotFound = errors.New("list item not found")
)

type ItemPostgres struct {
	db *gorm.DB
}

func NewTodoItemPostgres(db *gorm.DB) *ItemPostgres {
	return &ItemPostgres{
		db: db,
	}
}

func (ip *ItemPostgres) Create(item *domain.TodoItem) error {
	var list domain.TodoList
	if err := ip.db.Where(&domain.TodoList{Id: item.ListId}).First(&list).Error; err != nil {
		return ErrTodoListNotFound
	}

	if err := ip.db.Create(&item).Error; err != nil {
		return ErrCreateTodoItem
	}

	listItem := domain.ListsItem{
		ListId: item.ListId,
		ItemId: item.Id,
	}
	if err := ip.db.Create(&listItem).Error; err != nil {
		return ErrCreateListItem
	}

	return nil
}

func (ip *ItemPostgres) GetAll(listId int64) ([]*domain.TodoItem, error) {
	var listItems []domain.ListsItem
	if result := ip.db.Where(&domain.ListsItem{ListId: listId}).Find(&listItems); result.Error != nil {
		return nil, result.Error
	}

	var items []*domain.TodoItem
	for _, listItem := range listItems {
		var item *domain.TodoItem
		if result := ip.db.Where(&domain.TodoItem{Id: listItem.ItemId}).First(&item); result.Error == nil {
			items = append(items, item)
		}
	}

	return items, nil
}

func (ip *ItemPostgres) GetById(itemId int64) (*domain.TodoItem, int64, error) {
	var item domain.TodoItem
	if err := ip.db.Where(&domain.TodoItem{Id: itemId}).First(&item).Error; err != nil {
		return nil, 0, ErrTodoItemNotFound
	}

	var listItem *domain.ListsItem
	if err := ip.db.Where(&domain.ListsItem{ItemId: itemId}).First(&listItem).Error; err != nil {
		return nil, 0, ErrListItemNotFound
	}

	return &item, listItem.ListId, nil
}

func (ip *ItemPostgres) Delete(itemId int64) error {
	var item domain.TodoItem
	if result := ip.db.Where(&domain.TodoItem{Id: itemId}).First(&item); result.Error != nil {
		return result.Error
	}

	var listItem domain.ListsItem
	if result := ip.db.Where(&domain.ListsItem{ItemId: itemId}).First(&listItem); result.Error != nil {
		return result.Error
	}

	if err := ip.db.Delete(&listItem).Error; err != nil {
		return err
	}

	if err := ip.db.Delete(&item).Error; err != nil {
		return err
	}

	return nil
}

func (ip *ItemPostgres) Update(input *domain.TodoItem) error {
	if err := ip.db.Save(&input).Error; err != nil {
		return err
	}
	return nil
}
