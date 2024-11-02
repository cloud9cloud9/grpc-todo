package utils

import (
	"github.com/cloud9cloud9/go-grpc-todo/todo-svc/internal/domain"
	pb "github.com/cloud9cloud9/go-grpc-todo/todo-svc/internal/pb"
)

type Mapper interface {
	FromDomainLists(domainList []*domain.TodoList) []*pb.TodoList
	FromDomainItems(domainList []*domain.TodoItem) []*pb.TodoItem
}

type MapperImpl struct {
}

func NewMapper() Mapper {
	return &MapperImpl{}
}

func (m *MapperImpl) FromDomainLists(domainList []*domain.TodoList) []*pb.TodoList {
	var lists []*pb.TodoList
	for _, list := range domainList {
		lists = append(lists, &pb.TodoList{
			Id:    list.Id,
			Title: list.Title,
		})
	}

	return lists
}

func (m *MapperImpl) FromDomainItems(domainList []*domain.TodoItem) []*pb.TodoItem {
	var items []*pb.TodoItem
	for _, item := range domainList {
		items = append(items, &pb.TodoItem{
			Id:          item.Id,
			Title:       item.Title,
			Description: item.Description,
			Completed:   item.Done,
		})
	}

	return items
}
