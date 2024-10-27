package service

import (
	"context"
	"github.com/cloud9cloud9/go-grpc-todo/auth-svc/internal/domain"
	pb "github.com/cloud9cloud9/go-grpc-todo/auth-svc/internal/pb"
	"github.com/cloud9cloud9/go-grpc-todo/auth-svc/pkg/db"
	"github.com/cloud9cloud9/go-grpc-todo/auth-svc/pkg/utils"
	"net/http"
)

var (
	ErrUserNotFound = "User not found"
	ErrEmailExists  = "E-Mail already exists"
	ErrWrongPass    = "Wrong password"
)

type Server struct {
	Repo db.Repository
	pb.UnimplementedAuthServiceServer
}

func (s *Server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user domain.User

	if result := s.Repo.DB.Where(&domain.User{Email: in.Email}).First(&user); result.Error != nil {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  ErrUserNotFound,
		}, nil
	}

	if !utils.CompareHashAndPassword(user.Password, []byte(in.Password)) {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  ErrWrongPass,
		}, nil
	}

	token, err := utils.GenerateToken(user)
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
	var user domain.User

	if result := s.Repo.DB.Where(&domain.User{Email: in.Email}).First(&user); result.Error == nil {
		return &pb.RegisterResponse{
			Status: http.StatusConflict,
			Error:  ErrEmailExists,
		}, nil
	}

	user.Email = in.Email
	user.Password = utils.HashPassword(in.Password)

	s.Repo.DB.Create(&user)

	return &pb.RegisterResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *Server) Validate(ctx context.Context, in *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	var user domain.User

	claims, err := utils.ValidateToken(in.Token)
	if err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusUnauthorized,
			Error:  err.Error(),
		}, nil
	}

	if result := s.Repo.DB.Where(&domain.User{Id: claims.UserId}).First(&user); result.Error != nil {
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
