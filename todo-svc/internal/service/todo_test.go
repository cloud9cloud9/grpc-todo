package service

import (
	"context"
	"errors"
	"github.com/cloud9cloud9/go-grpc-todo/todo-svc/internal/domain"
	pb "github.com/cloud9cloud9/go-grpc-todo/todo-svc/internal/pb"
	mock_repository "github.com/cloud9cloud9/go-grpc-todo/todo-svc/internal/repository/mocks"
	"github.com/cloud9cloud9/go-grpc-todo/todo-svc/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestServer_CreateTodoList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	listRepo := mock_repository.NewMockTodoList(ctrl)

	serv := &Server{
		ListRepo: listRepo,
	}

	tests := []struct {
		name           string
		in             *pb.CreateTodoListRequest
		mockRepoSetup  func()
		expectedStatus int64
		expectedError  string
	}{
		{
			name: "Successful creation",
			in: &pb.CreateTodoListRequest{
				UserId: 1,
				Title:  "My Todo List",
			},
			mockRepoSetup: func() {
				listRepo.EXPECT().Create(int64(1), gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			name: "Failed creation",
			in: &pb.CreateTodoListRequest{
				UserId: 1,
				Title:  "My Todo List",
			},
			mockRepoSetup: func() {
				listRepo.EXPECT().Create(int64(1), gomock.Any()).Return(errors.New("error creating list"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "error creating list",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoSetup()
			resp, err := serv.CreateTodoList(context.Background(), tt.in)

			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedStatus, resp.Status)
				assert.NotNil(t, resp.List)
				assert.Equal(t, tt.in.Title, resp.List.Title)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedStatus, resp.Status)
				assert.Equal(t, tt.expectedError, resp.Error)
			}
		})
	}
}

func TestServer_GetTodoListById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	listRepo := mock_repository.NewMockTodoList(ctrl)

	serv := &Server{
		ListRepo: listRepo,
	}

	tests := []struct {
		name           string
		in             *pb.GetTodoListRequest
		mockRepoSetup  func()
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Success",
			in: &pb.GetTodoListRequest{
				Id:     1,
				UserId: 1,
			},
			mockRepoSetup: func() {
				listRepo.EXPECT().CheckUserAccessToList(int64(1), int64(1)).Return(nil)
				listRepo.EXPECT().GetById(int64(1)).Return(&domain.TodoList{Id: 1, Title: "My Todo List"}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name: "Not found",
			in: &pb.GetTodoListRequest{
				Id:     1,
				UserId: 1,
			},
			mockRepoSetup: func() {
				listRepo.EXPECT().GetById(int64(1)).Return(nil, errors.New("List not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "List not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoSetup()
			resp, err := serv.GetTodoListById(context.Background(), tt.in)
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tt.expectedStatus, int(resp.Status))
			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, resp.Error)
			}
		})
	}
}

func TestServer_GetTodoLists(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	listRepo := mock_repository.NewMockTodoList(ctrl)

	serv := &Server{
		ListRepo: listRepo,
		Mapper:   utils.NewMapper(),
	}

	tests := []struct {
		name           string
		in             *pb.GetTodoListsRequest
		mockRepoSetup  func()
		expectedStatus int
		expectedError  string
		expectedLists  []*pb.TodoList
	}{
		{
			name: "Success",
			in: &pb.GetTodoListsRequest{
				UserId: 1,
			},
			mockRepoSetup: func() {
				listRepo.EXPECT().GetAll(int64(1)).Return([]*domain.TodoList{
					{Id: 1, Title: "My First Todo"},
					{Id: 2, Title: "My Second Todo"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
			expectedLists: []*pb.TodoList{
				{Id: 1, Title: "My First Todo"},
				{Id: 2, Title: "My Second Todo"},
			},
		},
		{
			name: "Not found",
			in: &pb.GetTodoListsRequest{
				UserId: 1,
			},
			mockRepoSetup: func() {
				listRepo.EXPECT().GetAll(int64(1)).Return(nil, errors.New("List not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "List not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoSetup()
			resp, err := serv.GetTodoLists(context.Background(), tt.in)
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tt.expectedStatus, int(resp.Status))
			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, resp.Error)
			} else {
				assert.ElementsMatch(t, tt.expectedLists, resp.Lists)
			}
		})
	}
}

func TestServer_UpdateTodoList(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	listRepo := mock_repository.NewMockTodoList(ctrl)

	serv := &Server{
		ListRepo: listRepo,
	}

	tests := []struct {
		name           string
		in             *pb.UpdateTodoListRequest
		mockRepoSetup  func()
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Success",
			in: &pb.UpdateTodoListRequest{
				Id:     1,
				UserId: 1,
				Title:  "My Updated Todo",
			},
			mockRepoSetup: func() {
				listRepo.EXPECT().Update(int64(1), &domain.TodoList{Title: "My Updated Todo"}).Return(&domain.TodoList{Id: 1, Title: "My Updated Todo"}, nil)
				listRepo.EXPECT().CheckUserAccessToList(int64(1), int64(1)).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name: "Not found",
			in: &pb.UpdateTodoListRequest{
				Id:     1,
				UserId: 1,
				Title:  "My Updated Todo",
			},
			mockRepoSetup: func() {
				listRepo.EXPECT().Update(int64(1), &domain.TodoList{Title: "My Updated Todo"}).Return(nil, errors.New("List not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "List not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoSetup()
			resp, err := serv.UpdateTodoList(context.Background(), tt.in)
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tt.expectedStatus, int(resp.Status))
			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, resp.Error)
			}
		})
	}
}

func TestServer_DeleteTodoList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	listRepo := mock_repository.NewMockTodoList(ctrl)

	serv := &Server{
		ListRepo: listRepo,
	}

	tests := []struct {
		name           string
		in             *pb.DeleteTodoListRequest
		mockRepoSetup  func()
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Success",
			in: &pb.DeleteTodoListRequest{
				Id:     1,
				UserId: 1,
			},
			mockRepoSetup: func() {
				listRepo.EXPECT().Delete(int64(1)).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name: "Not found",
			in: &pb.DeleteTodoListRequest{
				Id:     1,
				UserId: 1,
			},
			mockRepoSetup: func() {
				listRepo.EXPECT().Delete(int64(1)).Return(errors.New("List not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "List not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoSetup()
			resp, err := serv.DeleteTodoList(context.Background(), tt.in)
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tt.expectedStatus, int(resp.Status))
			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, resp.Error)
			}
		})
	}
}

func TestServer_CreateTodoItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	itemRepo := mock_repository.NewMockTodoItem(ctrl)
	listRepo := mock_repository.NewMockTodoList(ctrl)
	serv := &Server{
		ItemRepo: itemRepo,
		ListRepo: listRepo,
	}

	tests := []struct {
		name           string
		in             *pb.CreateTodoItemRequest
		mockRepoSetup  func()
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Success",
			in: &pb.CreateTodoItemRequest{
				UserId:      1,
				ListId:      1,
				Title:       "My Todo Item",
				Description: "My Todo Item",
			},
			mockRepoSetup: func() {
				listRepo.EXPECT().CheckUserAccessToList(int64(1), int64(1)).Return(nil)
				itemRepo.EXPECT().Create(&domain.TodoItem{Title: "My Todo Item", Description: "My Todo Item", ListId: 1}).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name: "Not created",
			in: &pb.CreateTodoItemRequest{
				UserId:      1,
				ListId:      1,
				Title:       "My Todo Item",
				Description: "My Todo Item",
			},
			mockRepoSetup: func() {
				listRepo.EXPECT().CheckUserAccessToList(int64(1), int64(1)).Return(nil)
				itemRepo.EXPECT().Create(&domain.TodoItem{Title: "My Todo Item", Description: "My Todo Item", ListId: 1}).
					Return(errors.New("Item not created"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Item not created",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoSetup()
			resp, err := serv.CreateTodoItem(context.Background(), tt.in)
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tt.expectedStatus, int(resp.Status))
			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, resp.Error)
			}
		})
	}
}

func TestServer_GetTodoItemById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	itemRepo := mock_repository.NewMockTodoItem(ctrl)
	listRepo := mock_repository.NewMockTodoList(ctrl)
	serv := &Server{
		ItemRepo: itemRepo,
		ListRepo: listRepo,
	}

	tests := []struct {
		name           string
		in             *pb.GetTodoItemRequest
		mockRepoSetup  func()
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Success",
			in: &pb.GetTodoItemRequest{
				UserId: 1,
				Id:     1,
			},
			mockRepoSetup: func() {
				itemRepo.EXPECT().GetById(int64(1)).
					Return(
						&domain.TodoItem{Id: 1, Title: "My Todo Item", Description: "My Todo Item", ListId: 1},
						int64(1),
						nil,
					)
				listRepo.EXPECT().CheckUserAccessToList(int64(1), int64(1)).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name: "Not found",
			in: &pb.GetTodoItemRequest{
				UserId: 1,
				Id:     1,
			},
			mockRepoSetup: func() {
				itemRepo.EXPECT().GetById(int64(1)).
					Return(
						nil,
						int64(0),
						errors.New("Item not found"),
					)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Item not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoSetup()
			resp, err := serv.GetTodoItemById(context.Background(), tt.in)
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tt.expectedStatus, int(resp.Status))
			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, resp.Error)
			}
		})
	}
}

func TestServer_GetTodoItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	itemRepo := mock_repository.NewMockTodoItem(ctrl)
	listRepo := mock_repository.NewMockTodoList(ctrl)
	serv := &Server{
		ItemRepo: itemRepo,
		ListRepo: listRepo,
		Mapper:   utils.NewMapper(),
	}

	tests := []struct {
		name           string
		in             *pb.GetTodoItemsRequest
		mockRepoSetup  func()
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Success",
			in: &pb.GetTodoItemsRequest{
				UserId: 1,
				ListId: 1,
			},
			mockRepoSetup: func() {
				listRepo.EXPECT().CheckUserAccessToList(int64(1), int64(1)).Return(nil)
				itemRepo.EXPECT().GetAll(int64(1)).
					Return(
						[]*domain.TodoItem{{Id: 1, Title: "My Todo Item", Description: "My Todo Item", ListId: 1},
							{Id: 2, Title: "My Todo Item 2", Description: "My Todo Item 2", ListId: 1}},
						nil,
					)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name: "Not found",
			in: &pb.GetTodoItemsRequest{
				UserId: 1,
				ListId: 1,
			},
			mockRepoSetup: func() {
				listRepo.EXPECT().CheckUserAccessToList(int64(1), int64(1)).Return(nil)
				itemRepo.EXPECT().GetAll(int64(1)).
					Return(
						[]*domain.TodoItem{},
						errors.New("Item not found"),
					)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Item not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoSetup()
			resp, err := serv.GetTodoItems(context.Background(), tt.in)
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tt.expectedStatus, int(resp.Status))
			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, resp.Error)
			}
		})
	}
}

func TestServer_UpdateTodoItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	itemRepo := mock_repository.NewMockTodoItem(ctrl)
	listRepo := mock_repository.NewMockTodoList(ctrl)
	serv := &Server{
		ItemRepo: itemRepo,
		ListRepo: listRepo,
	}

	tests := []struct {
		name           string
		in             *pb.UpdateTodoItemRequest
		mockRepoSetup  func()
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Success",
			in: &pb.UpdateTodoItemRequest{
				UserId:      1,
				Id:          1,
				Title:       "My Todo Item",
				Completed:   true,
				Description: "My Todo Item",
			},
			mockRepoSetup: func() {
				itemRepo.EXPECT().GetById(int64(1)).
					Return(
						&domain.TodoItem{Id: 1, Title: "My Todo Item", Description: "My Todo Item", ListId: 1},
						int64(1),
						nil,
					)
				listRepo.EXPECT().CheckUserAccessToList(int64(1), int64(1)).Return(nil)
				itemRepo.EXPECT().Update(&domain.TodoItem{Id: 1, Title: "My Todo Item", Description: "My Todo Item", ListId: 1, Done: true}).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name: "Not found",
			in: &pb.UpdateTodoItemRequest{
				UserId:      1,
				Id:          1,
				Title:       "My Todo Item",
				Completed:   true,
				Description: "My Todo Item",
			},
			mockRepoSetup: func() {
				itemRepo.EXPECT().GetById(int64(1)).
					Return(
						nil,
						int64(0),
						errors.New("Item not found"),
					)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Item not found",
		},
		{
			name: "error updating item",
			in: &pb.UpdateTodoItemRequest{
				UserId:      1,
				Id:          1,
				Title:       "My Todo Item",
				Completed:   true,
				Description: "My Todo Item",
			},
			mockRepoSetup: func() {
				itemRepo.EXPECT().GetById(int64(1)).
					Return(
						&domain.TodoItem{Id: 1, Title: "My Todo Item", Description: "My Todo Item", ListId: 1},
						int64(1),
						nil,
					)
				listRepo.EXPECT().CheckUserAccessToList(int64(1), int64(1)).Return(nil)
				itemRepo.EXPECT().Update(&domain.TodoItem{Id: 1, Title: "My Todo Item", Description: "My Todo Item", ListId: 1, Done: true}).
					Return(errors.New("error updating item"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "error updating item",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoSetup()
			resp, err := serv.UpdateTodoItem(context.Background(), tt.in)
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tt.expectedStatus, int(resp.Status))
			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, resp.Error)
			}
		})
	}
}

func TestServer_DeleteTodoItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	itemRepo := mock_repository.NewMockTodoItem(ctrl)
	listRepo := mock_repository.NewMockTodoList(ctrl)
	serv := &Server{
		ItemRepo: itemRepo,
		ListRepo: listRepo,
		Mapper:   utils.NewMapper(),
	}

	tests := []struct {
		name           string
		in             *pb.DeleteTodoItemRequest
		mockRepoSetup  func()
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Success",
			in: &pb.DeleteTodoItemRequest{
				UserId: 1,
				Id:     1,
			},
			mockRepoSetup: func() {
				itemRepo.EXPECT().GetById(int64(1)).
					Return(
						&domain.TodoItem{Id: 1, Title: "My Todo Item", Description: "My Todo Item", ListId: 1},
						int64(1),
						nil,
					)
				listRepo.EXPECT().CheckUserAccessToList(int64(1), int64(1)).Return(nil)
				itemRepo.EXPECT().Delete(int64(1)).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name: "Not found",
			in: &pb.DeleteTodoItemRequest{
				UserId: 1,
				Id:     1,
			},
			mockRepoSetup: func() {
				itemRepo.EXPECT().GetById(int64(1)).
					Return(
						nil,
						int64(0),
						errors.New("Item not found"),
					)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Item not found",
		},
		{
			name: "error deleting item",
			in: &pb.DeleteTodoItemRequest{
				UserId: 1,
				Id:     1,
			},
			mockRepoSetup: func() {
				itemRepo.EXPECT().GetById(int64(1)).
					Return(
						&domain.TodoItem{Id: 1, Title: "My Todo Item", Description: "My Todo Item", ListId: 1},
						int64(1),
						nil,
					)
				listRepo.EXPECT().CheckUserAccessToList(int64(1), int64(1)).Return(nil)
				itemRepo.EXPECT().Delete(int64(1)).Return(errors.New("error deleting item"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "error deleting item",
		},
		{
			name: "error accessing list",
			in: &pb.DeleteTodoItemRequest{
				UserId: 1,
				Id:     1,
			},
			mockRepoSetup: func() {
				itemRepo.EXPECT().GetById(int64(1)).
					Return(
						&domain.TodoItem{Id: 1, Title: "My Todo Item", Description: "My Todo Item", ListId: 1},
						int64(1),
						nil,
					)
				listRepo.EXPECT().CheckUserAccessToList(int64(1), int64(1)).Return(errors.New("User does not have access to this list"))
			},
			expectedStatus: http.StatusForbidden,
			expectedError:  "User does not have access to this list",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoSetup()
			resp, err := serv.DeleteTodoItem(context.Background(), tt.in)
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tt.expectedStatus, int(resp.Status))
			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, resp.Error)
			}
		})
	}
}
