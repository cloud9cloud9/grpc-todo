package mocks

import (
	"context"
	pb "github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/todo/pb"
	"google.golang.org/grpc"
)

type MockTodoServiceClient struct {
	CreateTodoItemFunc  func(ctx context.Context, in *pb.CreateTodoItemRequest) (*pb.CreateTodoItemResponse, error)
	CreateTodoListFunc  func(ctx context.Context, in *pb.CreateTodoListRequest) (*pb.CreateTodoListResponse, error)
	GetTodoListByIdFunc func(ctx context.Context, in *pb.GetTodoListRequest) (*pb.GetTodoListResponse, error)
	GetTodoListsFunc    func(ctx context.Context, in *pb.GetTodoListsRequest) (*pb.GetTodoListsResponse, error)
	UpdateTodoListFunc  func(ctx context.Context, in *pb.UpdateTodoListRequest) (*pb.UpdateTodoListResponse, error)
	DeleteTodoListFunc  func(ctx context.Context, in *pb.DeleteTodoListRequest) (*pb.DeleteTodoListResponse, error)
	DeleteTodoItemFunc  func(ctx context.Context, in *pb.DeleteTodoItemRequest) (*pb.DeleteTodoItemResponse, error)
	UpdateTodoItemFunc  func(ctx context.Context, in *pb.UpdateTodoItemRequest) (*pb.UpdateTodoItemResponse, error)
	GetTodoItemByIdFunc func(ctx context.Context, in *pb.GetTodoItemRequest) (*pb.GetTodoItemResponse, error)
	GetTodoItemsFunc    func(ctx context.Context, in *pb.GetTodoItemsRequest) (*pb.GetTodoItemsResponse, error)
}

func (m *MockTodoServiceClient) CreateTodoItem(ctx context.Context, in *pb.CreateTodoItemRequest, opts ...grpc.CallOption) (*pb.CreateTodoItemResponse, error) {
	return m.CreateTodoItemFunc(ctx, in)
}
func (m *MockTodoServiceClient) CreateTodoList(ctx context.Context, in *pb.CreateTodoListRequest, opts ...grpc.CallOption) (*pb.CreateTodoListResponse, error) {
	return m.CreateTodoListFunc(ctx, in)
}
func (m *MockTodoServiceClient) GetTodoListById(ctx context.Context, in *pb.GetTodoListRequest, opts ...grpc.CallOption) (*pb.GetTodoListResponse, error) {
	return m.GetTodoListByIdFunc(ctx, in)
}
func (m *MockTodoServiceClient) GetTodoLists(ctx context.Context, in *pb.GetTodoListsRequest, opts ...grpc.CallOption) (*pb.GetTodoListsResponse, error) {
	return m.GetTodoListsFunc(ctx, in)
}
func (m *MockTodoServiceClient) UpdateTodoList(ctx context.Context, in *pb.UpdateTodoListRequest, opts ...grpc.CallOption) (*pb.UpdateTodoListResponse, error) {
	return m.UpdateTodoListFunc(ctx, in)
}
func (m *MockTodoServiceClient) DeleteTodoList(ctx context.Context, in *pb.DeleteTodoListRequest, opts ...grpc.CallOption) (*pb.DeleteTodoListResponse, error) {
	return m.DeleteTodoListFunc(ctx, in)
}
func (m *MockTodoServiceClient) DeleteTodoItem(ctx context.Context, in *pb.DeleteTodoItemRequest, opts ...grpc.CallOption) (*pb.DeleteTodoItemResponse, error) {
	return m.DeleteTodoItemFunc(ctx, in)
}
func (m *MockTodoServiceClient) UpdateTodoItem(ctx context.Context, in *pb.UpdateTodoItemRequest, opts ...grpc.CallOption) (*pb.UpdateTodoItemResponse, error) {
	return m.UpdateTodoItemFunc(ctx, in)
}
func (m *MockTodoServiceClient) GetTodoItemById(ctx context.Context, in *pb.GetTodoItemRequest, opts ...grpc.CallOption) (*pb.GetTodoItemResponse, error) {
	return m.GetTodoItemByIdFunc(ctx, in)
}
func (m *MockTodoServiceClient) GetTodoItems(ctx context.Context, in *pb.GetTodoItemsRequest, opts ...grpc.CallOption) (*pb.GetTodoItemsResponse, error) {
	return m.GetTodoItemsFunc(ctx, in)
}
