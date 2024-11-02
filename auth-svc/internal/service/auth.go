package service

import (
	"context"
	"github.com/cloud9cloud9/go-grpc-todo/auth-svc/internal/domain"
	pb "github.com/cloud9cloud9/go-grpc-todo/auth-svc/internal/pb"
	"github.com/cloud9cloud9/go-grpc-todo/auth-svc/internal/repository"
	"github.com/cloud9cloud9/go-grpc-todo/auth-svc/internal/security"
	"net/http"
)

var (
	ErrUserNotFound = "User not found"
	ErrEmailExists  = "E-Mail already exists"
	ErrWrongPass    = "Wrong password"
)

type Server struct {
	Repo       repository.UserRepository
	AuthHelper security.AuthHelper
	pb.UnimplementedAuthServiceServer
}

func (s *Server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.Repo.FindByEmail(in.Email)
	if err != nil {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  ErrUserNotFound,
		}, nil
	}

	if !s.AuthHelper.CompareHashAndPassword(user.Password, []byte(in.Password)) {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  ErrWrongPass,
		}, nil
	}

	token, err := s.AuthHelper.GenerateToken(user)
	if err != nil {
		return &pb.LoginResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	return &pb.LoginResponse{
		Status: http.StatusOK,
		Token:  token,
	}, nil
}

func (s *Server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user, err := s.Repo.FindByEmail(in.Email)
	if err == nil && user != nil {
		return &pb.RegisterResponse{
			Status: http.StatusConflict,
			Error:  ErrEmailExists,
		}, nil
	}

	user = &domain.User{
		Email:    in.Email,
		Password: s.AuthHelper.HashPassword(in.Password),
	}

	err = s.Repo.CreateUser(user)
	if err != nil {
		return &pb.RegisterResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	return &pb.RegisterResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *Server) Validate(ctx context.Context, in *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	claims, err := s.AuthHelper.ValidateToken(in.Token)
	if err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusUnauthorized,
			Error:  err.Error(),
		}, nil
	}

	user, err := s.Repo.FindByID(claims.UserId)
	if err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusNotFound,
			Error:  ErrUserNotFound,
		}, nil
	}

	return &pb.ValidateResponse{
		Status: http.StatusOK,
		UserId: user.Id,
	}, nil
}
