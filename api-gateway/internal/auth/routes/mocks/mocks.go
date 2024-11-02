package mocks

import (
	"context"
	pb "github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/auth/pb"
	"google.golang.org/grpc"
)

type MockAuthServiceClient struct {
	LoginFunc    func(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	RegisterFunc func(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error)
	ValidateFunc func(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error)
}

func (m *MockAuthServiceClient) Login(ctx context.Context, in *pb.LoginRequest, opts ...grpc.CallOption) (*pb.LoginResponse, error) {
	return m.LoginFunc(ctx, in)
}

func (m *MockAuthServiceClient) Register(ctx context.Context, in *pb.RegisterRequest, opts ...grpc.CallOption) (*pb.RegisterResponse, error) {
	return m.RegisterFunc(ctx, in)
}

func (m *MockAuthServiceClient) Validate(ctx context.Context, in *pb.ValidateRequest, opts ...grpc.CallOption) (*pb.ValidateResponse, error) {
	return m.ValidateFunc(ctx, in)
}
