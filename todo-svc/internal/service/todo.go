package service

import (
	"context"
	"github.com/cloud9cloud9/go-grpc-todo/todo-svc/internal/domain"
	pb "github.com/cloud9cloud9/go-grpc-todo/todo-svc/internal/pb"
	"github.com/cloud9cloud9/go-grpc-todo/todo-svc/pkg/db"
	"net/http"
)

var (
	errForbidden    = "User does not have access to this list"
	errListNotFound = "List not found"
	errItemNotFound = "Item not found"
)

type Server struct {
	Repo db.Repository
	pb.UnimplementedTodoServiceServer
}

func (s *Server) CreateTodoList(ctx context.Context, in *pb.CreateTodoListRequest) (*pb.CreateTodoListResponse, error) {
	var list domain.TodoList
	list.Title = in.Title

	tx := s.Repo.DB.Begin()
	if err := tx.Create(&list).Error; err != nil {
		tx.Rollback()
		return &pb.CreateTodoListResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	userList := domain.UsersList{
		UserId: in.UserId,
		ListId: list.Id,
	}

	if err := tx.Create(&userList).Error; err != nil {
		tx.Rollback()
		return &pb.CreateTodoListResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	if err := tx.Commit().Error; err != nil {
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
	var list domain.TodoList

	if result := s.Repo.DB.Where(&domain.TodoList{Id: in.Id}).First(&list); result.Error != nil {
		return &pb.GetTodoListResponse{
			Status: http.StatusNotFound,
			Error:  result.Error.Error(),
		}, nil
	}

	var userList domain.UsersList
	if result := s.Repo.DB.Where(&domain.UsersList{UserId: in.UserId, ListId: list.Id}).First(&userList); result.Error != nil {
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
	var userLists []domain.UsersList
	if result := s.Repo.DB.Where(&domain.UsersList{UserId: in.UserId}).Find(&userLists); result.Error != nil {
		return &pb.GetTodoListsResponse{
			Status: http.StatusNotFound,
			Error:  errListNotFound,
		}, nil
	}

	var lists []*pb.TodoList

	for _, userList := range userLists {
		var list domain.TodoList
		if result := s.Repo.DB.Where(&domain.TodoList{Id: userList.ListId}).First(&list); result.Error == nil {
			lists = append(lists, &pb.TodoList{
				Id:    list.Id,
				Title: list.Title,
			})
		}
	}

	return &pb.GetTodoListsResponse{
		Lists:  lists,
		Status: http.StatusOK,
	}, nil
}

func (s *Server) UpdateTodoList(ctx context.Context, in *pb.UpdateTodoListRequest) (*pb.UpdateTodoListResponse, error) {
	var list domain.TodoList
	if result := s.Repo.DB.Where(&domain.TodoList{Id: in.Id}).First(&list); result.Error != nil {
		return &pb.UpdateTodoListResponse{
			Status: http.StatusNotFound,
			Error:  errListNotFound,
		}, nil
	}

	var userList domain.UsersList
	if result := s.Repo.DB.Where(&domain.UsersList{UserId: in.UserId, ListId: in.Id}).First(&userList); result.Error != nil {
		return &pb.UpdateTodoListResponse{
			Status: http.StatusForbidden,
			Error:  errForbidden,
		}, nil
	}

	list.Title = in.Title

	if err := s.Repo.DB.Save(&list).Error; err != nil {
		return &pb.UpdateTodoListResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
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
	var list domain.TodoList
	if result := s.Repo.DB.Where(&domain.TodoList{Id: in.Id}).First(&list); result.Error != nil {
		return &pb.DeleteTodoListResponse{
			Status: http.StatusNotFound,
			Error:  errListNotFound,
		}, nil
	}

	var userList domain.UsersList
	if result := s.Repo.DB.Where(&domain.UsersList{UserId: in.UserId, ListId: in.Id}).First(&userList); result.Error != nil {
		return &pb.DeleteTodoListResponse{
			Status: http.StatusForbidden,
			Error:  errForbidden,
		}, nil
	}

	if err := s.Repo.DB.Delete(&list).Error; err != nil {
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
	var list domain.TodoList

	if result := s.Repo.DB.Where(&domain.TodoList{Id: in.ListId}).First(&list); result.Error != nil {
		return &pb.CreateTodoItemResponse{
			Item:   nil,
			Status: http.StatusNotFound,
			Error:  errListNotFound,
		}, nil
	}

	var userList domain.UsersList
	if result := s.Repo.DB.Where(&domain.UsersList{UserId: in.UserId, ListId: in.ListId}).First(&userList); result.Error != nil {
		return &pb.CreateTodoItemResponse{
			Item:   nil,
			Status: http.StatusForbidden,
			Error:  errForbidden,
		}, nil
	}

	item := domain.TodoItem{
		Title:       in.Title,
		Description: in.Description,
		Done:        false,
		ListId:      list.Id,
	}

	if err := s.Repo.DB.Create(&item).Error; err != nil {
		return &pb.CreateTodoItemResponse{
			Item:   nil,
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	listItem := domain.ListsItem{
		ListId: list.Id,
		ItemId: item.Id,
	}

	if err := s.Repo.DB.Create(&listItem).Error; err != nil {
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
	var item domain.TodoItem

	if result := s.Repo.DB.Where(&domain.TodoItem{Id: in.Id}).First(&item); result.Error != nil {
		return &pb.GetTodoItemResponse{
			Item:   nil,
			Status: http.StatusNotFound,
			Error:  errItemNotFound,
		}, nil
	}

	var listItem domain.ListsItem
	if result := s.Repo.DB.Where(&domain.ListsItem{ItemId: in.Id}).Find(&listItem); result.Error != nil {
		return &pb.GetTodoItemResponse{
			Item:   nil,
			Status: http.StatusForbidden,
			Error:  errItemNotFound,
		}, nil
	}

	var userHasAccess bool
	var userList domain.UsersList
	if result := s.Repo.DB.Where(&domain.UsersList{UserId: in.UserId, ListId: listItem.ListId}).First(&userList); result.Error == nil {
		userHasAccess = true
	}

	if !userHasAccess {
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
	var userList domain.UsersList

	if result := s.Repo.DB.Where(&domain.UsersList{UserId: in.UserId, ListId: in.ListId}).First(&userList); result.Error != nil {
		return &pb.GetTodoItemsResponse{
			Items:  nil,
			Status: http.StatusForbidden,
			Error:  errForbidden,
		}, nil
	}

	var listItems []domain.ListsItem
	if result := s.Repo.DB.Where(&domain.ListsItem{ListId: in.ListId}).Find(&listItems); result.Error != nil {
		return &pb.GetTodoItemsResponse{
			Items:  nil,
			Status: http.StatusNotFound,
			Error:  errItemNotFound,
		}, nil
	}

	var items []*pb.TodoItem

	for _, listItem := range listItems {
		var item domain.TodoItem
		if result := s.Repo.DB.Where(&domain.TodoItem{Id: listItem.ItemId}).First(&item); result.Error == nil {
			items = append(items, &pb.TodoItem{
				Id:          item.Id,
				Title:       item.Title,
				Description: item.Description,
			})
		}
	}

	return &pb.GetTodoItemsResponse{
		Items:  items,
		Status: http.StatusOK,
	}, nil
}

func (s *Server) UpdateTodoItem(ctx context.Context, in *pb.UpdateTodoItemRequest) (*pb.UpdateTodoItemResponse, error) {
	var item domain.TodoItem

	if result := s.Repo.DB.Where(&domain.TodoItem{Id: in.Id}).First(&item); result.Error != nil {
		return &pb.UpdateTodoItemResponse{
			Item:   nil,
			Status: http.StatusNotFound,
			Error:  errItemNotFound,
		}, nil
	}

	var listItem domain.ListsItem
	if result := s.Repo.DB.Where(&domain.ListsItem{ItemId: item.Id}).First(&listItem); result.Error != nil {
		return &pb.UpdateTodoItemResponse{
			Item:   nil,
			Status: http.StatusForbidden,
			Error:  errItemNotFound,
		}, nil
	}

	var userHasAccess bool
	var userList domain.UsersList
	if result := s.Repo.DB.Where(&domain.UsersList{UserId: in.UserId, ListId: listItem.ListId}).First(&userList); result.Error == nil {
		userHasAccess = true
	}

	if !userHasAccess {
		return &pb.UpdateTodoItemResponse{
			Item:   nil,
			Status: http.StatusForbidden,
			Error:  errForbidden,
		}, nil
	}

	item.Title = in.Title
	item.Description = in.Description
	item.Done = in.Completed

	if err := s.Repo.DB.Save(&item).Error; err != nil {
		return &pb.UpdateTodoItemResponse{
			Item:   nil,
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
	var item domain.TodoItem

	if result := s.Repo.DB.Where(&domain.TodoItem{Id: in.Id}).First(&item); result.Error != nil {
		return &pb.DeleteTodoItemResponse{
			Success: false,
			Status:  http.StatusNotFound,
			Error:   errItemNotFound,
		}, nil
	}

	var listItem domain.ListsItem
	if result := s.Repo.DB.Where(&domain.ListsItem{ItemId: item.Id}).First(&listItem); result.Error != nil {
		return &pb.DeleteTodoItemResponse{
			Success: false,
			Status:  http.StatusForbidden,
			Error:   errItemNotFound,
		}, nil
	}

	var userList domain.UsersList
	if result := s.Repo.DB.Where(&domain.UsersList{UserId: in.UserId, ListId: listItem.ListId}).First(&userList); result.Error != nil {
		return &pb.DeleteTodoItemResponse{
			Success: false,
			Status:  http.StatusForbidden,
			Error:   errForbidden,
		}, nil
	}

	if err := s.Repo.DB.Delete(&listItem).Error; err != nil {
		return &pb.DeleteTodoItemResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   err.Error(),
		}, nil
	}

	if err := s.Repo.DB.Delete(&item).Error; err != nil {
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
