package service

import (
	"context"
	"github.com/cloud9cloud9/go-grpc-todo/todo-svc/internal/domain"
	pb "github.com/cloud9cloud9/go-grpc-todo/todo-svc/internal/pb"
	"github.com/cloud9cloud9/go-grpc-todo/todo-svc/internal/repository"
	"github.com/cloud9cloud9/go-grpc-todo/todo-svc/pkg/utils"
	"net/http"
)

var (
	errForbidden    = "User does not have access to this list"
	errListNotFound = "List not found"
	errItemNotFound = "Item not found"
)

type Server struct {
	ListRepo repository.TodoList
	ItemRepo repository.TodoItem
	pb.UnimplementedTodoServiceServer
	Mapper utils.Mapper
}

func (s *Server) CreateTodoList(ctx context.Context, in *pb.CreateTodoListRequest) (*pb.CreateTodoListResponse, error) {
	var list domain.TodoList
	list.Title = in.Title

	if err := s.ListRepo.Create(in.UserId, &list); err != nil {
		return &pb.CreateTodoListResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	return &pb.CreateTodoListResponse{
		List: &pb.TodoList{
			Id:    list.Id,
			Title: list.Title,
		},
		Status: http.StatusCreated,
	}, nil
}

func (s *Server) GetTodoListById(ctx context.Context, in *pb.GetTodoListRequest) (*pb.GetTodoListResponse, error) {
	list, err := s.ListRepo.GetById(in.Id)
	if err != nil {
		return &pb.GetTodoListResponse{
			Status: http.StatusNotFound,
			Error:  errListNotFound,
		}, nil
	}

	if s.ListRepo.CheckUserAccessToList(in.UserId, list.Id) != nil {
		return &pb.GetTodoListResponse{
			Status: http.StatusForbidden,
			Error:  errForbidden,
		}, nil
	}

	return &pb.GetTodoListResponse{
		List: &pb.TodoList{
			Id:    list.Id,
			Title: list.Title,
		},
		Status: http.StatusOK,
	}, nil
}

func (s *Server) GetTodoLists(ctx context.Context, in *pb.GetTodoListsRequest) (*pb.GetTodoListsResponse, error) {
	todoLists, err := s.ListRepo.GetAll(in.UserId)
	if err != nil {
		return &pb.GetTodoListsResponse{
			Status: http.StatusNotFound,
			Error:  errListNotFound,
		}, nil
	}

	lists := s.Mapper.FromDomainLists(todoLists)

	return &pb.GetTodoListsResponse{
		Lists:  lists,
		Status: http.StatusOK,
	}, nil
}

func (s *Server) UpdateTodoList(ctx context.Context, in *pb.UpdateTodoListRequest) (*pb.UpdateTodoListResponse, error) {
	list, err := s.ListRepo.Update(in.Id, &domain.TodoList{
		Title: in.Title,
	})
	if err != nil {
		return &pb.UpdateTodoListResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	err = s.ListRepo.CheckUserAccessToList(in.UserId, in.Id)
	if err != nil {
		return &pb.UpdateTodoListResponse{
			Status: http.StatusForbidden,
			Error:  errForbidden,
		}, nil
	}

	return &pb.UpdateTodoListResponse{
		List: &pb.TodoList{
			Id:    list.Id,
			Title: list.Title,
		},
		Status: http.StatusOK,
	}, nil
}

func (s *Server) DeleteTodoList(ctx context.Context, in *pb.DeleteTodoListRequest) (*pb.DeleteTodoListResponse, error) {
	if err := s.ListRepo.Delete(in.Id); err != nil {
		return &pb.DeleteTodoListResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	return &pb.DeleteTodoListResponse{
		Status:  http.StatusOK,
		Success: true,
	}, nil
}

func (s *Server) CreateTodoItem(ctx context.Context, in *pb.CreateTodoItemRequest) (*pb.CreateTodoItemResponse, error) {
	if err := s.ListRepo.CheckUserAccessToList(in.UserId, in.ListId); err != nil {
		return &pb.CreateTodoItemResponse{
			Item:   nil,
			Status: http.StatusForbidden,
			Error:  err.Error(),
		}, nil
	}

	item := &domain.TodoItem{
		Title:       in.Title,
		Description: in.Description,
		Done:        false,
		ListId:      in.ListId,
	}

	if err := s.ItemRepo.Create(item); err != nil {
		return &pb.CreateTodoItemResponse{
			Item:   nil,
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	return &pb.CreateTodoItemResponse{
		Item: &pb.TodoItem{
			Id:          item.Id,
			Title:       item.Title,
			Description: item.Description,
		},
		Status: http.StatusOK,
	}, nil
}

func (s *Server) GetTodoItemById(ctx context.Context, in *pb.GetTodoItemRequest) (*pb.GetTodoItemResponse, error) {
	item, _, err := s.ItemRepo.GetById(in.Id)
	if err != nil {
		return &pb.GetTodoItemResponse{
			Item:   nil,
			Status: http.StatusNotFound,
			Error:  errItemNotFound,
		}, nil
	}

	if err := s.ListRepo.CheckUserAccessToList(in.UserId, item.ListId); err != nil {
		return &pb.GetTodoItemResponse{
			Item:   nil,
			Status: http.StatusForbidden,
			Error:  errForbidden,
		}, nil
	}

	return &pb.GetTodoItemResponse{
		Item: &pb.TodoItem{
			Id:          item.Id,
			Title:       item.Title,
			Description: item.Description,
		},
		Status: http.StatusOK,
	}, nil
}

func (s *Server) GetTodoItems(ctx context.Context, in *pb.GetTodoItemsRequest) (*pb.GetTodoItemsResponse, error) {
	if err := s.ListRepo.CheckUserAccessToList(in.UserId, in.ListId); err != nil {
		return &pb.GetTodoItemsResponse{
			Status: http.StatusForbidden,
			Error:  errForbidden,
		}, nil
	}

	listItems, err := s.ItemRepo.GetAll(in.ListId)
	if err != nil {
		return &pb.GetTodoItemsResponse{
			Status: http.StatusInternalServerError,
			Error:  errItemNotFound,
		}, nil
	}

	items := s.Mapper.FromDomainItems(listItems)

	return &pb.GetTodoItemsResponse{
		Items:  items,
		Status: http.StatusOK,
	}, nil
}

func (s *Server) UpdateTodoItem(ctx context.Context, in *pb.UpdateTodoItemRequest) (*pb.UpdateTodoItemResponse, error) {
	item, listId, err := s.ItemRepo.GetById(in.Id)
	if err != nil {
		return &pb.UpdateTodoItemResponse{
			Status: http.StatusNotFound,
			Error:  errItemNotFound,
		}, nil
	}

	if err := s.ListRepo.CheckUserAccessToList(in.UserId, listId); err != nil {
		return &pb.UpdateTodoItemResponse{
			Status: http.StatusForbidden,
			Error:  errForbidden,
		}, nil
	}

	item.Id = in.Id
	item.Title = in.Title
	item.Description = in.Description
	item.Done = in.Completed

	if err := s.ItemRepo.Update(item); err != nil {
		return &pb.UpdateTodoItemResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	return &pb.UpdateTodoItemResponse{
		Item: &pb.TodoItem{
			Id:          item.Id,
			Title:       item.Title,
			Description: item.Description,
		},
		Status: http.StatusOK,
	}, nil
}

func (s *Server) DeleteTodoItem(ctx context.Context, in *pb.DeleteTodoItemRequest) (*pb.DeleteTodoItemResponse, error) {
	item, listId, err := s.ItemRepo.GetById(in.Id)
	if err != nil {
		return &pb.DeleteTodoItemResponse{
			Success: false,
			Status:  http.StatusNotFound,
			Error:   errItemNotFound,
		}, nil
	}

	if err := s.ListRepo.CheckUserAccessToList(in.UserId, listId); err != nil {
		return &pb.DeleteTodoItemResponse{
			Success: false,
			Status:  http.StatusForbidden,
			Error:   errForbidden,
		}, nil
	}

	if err := s.ItemRepo.Delete(item.Id); err != nil {
		return &pb.DeleteTodoItemResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   err.Error(),
		}, nil
	}

	return &pb.DeleteTodoItemResponse{
		Success: true,
		Status:  http.StatusOK,
	}, nil
}
