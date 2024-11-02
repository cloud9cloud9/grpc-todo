package service

import (
	"context"
	"errors"
	"github.com/cloud9cloud9/go-grpc-todo/auth-svc/internal/domain"
	pb "github.com/cloud9cloud9/go-grpc-todo/auth-svc/internal/pb"
	mock "github.com/cloud9cloud9/go-grpc-todo/auth-svc/internal/repository/mocks"
	"github.com/cloud9cloud9/go-grpc-todo/auth-svc/internal/security"
	mock_utils "github.com/cloud9cloud9/go-grpc-todo/auth-svc/internal/security/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const (
	testEmail      = "test@example.com"
	testPassword   = "password"
	hashedPassword = "hashed_password"
	validToken     = "valid_token"
	invalidToken   = "invalid_token"
)

func TestServer_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockUserRepository(ctrl)
	mockAuthHelper := mock_utils.NewMockAuthHelper(ctrl)
	s := Server{
		Repo:       mockRepo,
		AuthHelper: mockAuthHelper,
	}

	tests := []struct {
		name           string
		email          string
		password       string
		mockRepoSetup  func()
		mockAuthSetup  func()
		expectedStatus int
		expectedError  string
	}{
		{
			name:     "Success",
			email:    testEmail,
			password: testPassword,
			mockRepoSetup: func() {
				user := &domain.User{Email: testEmail, Password: hashedPassword}
				mockRepo.EXPECT().FindByEmail(testEmail).Return(user, nil)
			},
			mockAuthSetup: func() {
				mockAuthHelper.EXPECT().CompareHashAndPassword(hashedPassword, []byte(testPassword)).Return(true)
				mockAuthHelper.EXPECT().GenerateToken(gomock.Any()).Return("valid_token", nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:     "User Not Found",
			email:    testEmail,
			password: testPassword,
			mockRepoSetup: func() {
				user := &domain.User{}
				mockRepo.EXPECT().FindByEmail(testEmail).Return(user, errors.New("user not found"))
			},
			mockAuthSetup:  func() {},
			expectedStatus: http.StatusNotFound,
			expectedError:  ErrUserNotFound,
		},
		{
			name:     "Wrong Password",
			email:    "test@example.com",
			password: "wrongpassword",
			mockRepoSetup: func() {
				user := &domain.User{Email: testEmail, Password: hashedPassword}
				mockRepo.EXPECT().FindByEmail(testEmail).Return(user, nil)
			},
			mockAuthSetup: func() {
				mockAuthHelper.EXPECT().CompareHashAndPassword(hashedPassword, []byte("wrongpassword")).Return(false)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  ErrWrongPass,
		},
		{
			name:     "Token Generation Error",
			email:    "test@example.com",
			password: "password123",
			mockRepoSetup: func() {
				user := &domain.User{Email: testEmail, Password: hashedPassword}
				mockRepo.EXPECT().FindByEmail(testEmail).Return(user, nil)
			},
			mockAuthSetup: func() {
				mockAuthHelper.EXPECT().CompareHashAndPassword(gomock.Any(), gomock.Any()).Return(true)
				mockAuthHelper.EXPECT().GenerateToken(gomock.Any()).Return("", errors.New("token generation error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "token generation error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoSetup()
			tt.mockAuthSetup()

			req := &pb.LoginRequest{Email: tt.email, Password: tt.password}
			resp, err := s.Login(context.Background(), req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, int(resp.Status))
			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, resp.Error)
			} else {
				assert.Equal(t, "valid_token", resp.Token)
			}
		})
	}
}

func TestServer_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockUserRepository(ctrl)
	mockAuthHelper := mock_utils.NewMockAuthHelper(ctrl)
	s := Server{
		Repo:       mockRepo,
		AuthHelper: mockAuthHelper,
	}

	tests := []struct {
		name           string
		email          string
		password       string
		mockRepoSetup  func()
		mockAuthSetup  func()
		expectedStatus int
		expectedError  string
	}{
		{
			name:     "Success",
			email:    testEmail,
			password: testPassword,
			mockRepoSetup: func() {
				mockRepo.EXPECT().FindByEmail(testEmail).Return(nil, errors.New("user not found"))
				mockRepo.EXPECT().CreateUser(gomock.Any()).Return(nil)
			},
			mockAuthSetup: func() {
				mockAuthHelper.EXPECT().HashPassword(testPassword).Return(hashedPassword)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			name:     "User Already Exists",
			email:    testEmail,
			password: testPassword,
			mockRepoSetup: func() {
				user := &domain.User{Email: testEmail}
				mockRepo.EXPECT().FindByEmail(testEmail).Return(user, nil)
			},
			mockAuthSetup:  func() {},
			expectedStatus: http.StatusConflict,
			expectedError:  ErrEmailExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoSetup()
			tt.mockAuthSetup()

			req := &pb.RegisterRequest{Email: tt.email, Password: tt.password}
			resp, err := s.Register(context.Background(), req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, int(resp.Status))
			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, resp.Error)
			}
		})
	}
}

func TestServer_Validate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockUserRepository(ctrl)
	mockAuthHelper := mock_utils.NewMockAuthHelper(ctrl)
	s := Server{
		Repo:       mockRepo,
		AuthHelper: mockAuthHelper,
	}

	tests := []struct {
		name           string
		token          string
		mockAuthSetup  func()
		mockRepoSetup  func()
		expectedStatus int
		expectedError  string
	}{
		{
			name:  "Success",
			token: "valid_token",
			mockAuthSetup: func() {
				mockAuthHelper.EXPECT().ValidateToken(validToken).Return(&security.TokenClaims{UserId: 1, Email: testEmail}, nil)
			},
			mockRepoSetup: func() {
				id := int64(1)
				mockRepo.EXPECT().FindByID(id).Return(&domain.User{Id: 1, Email: testEmail}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:  "Invalid Token",
			token: "invalid_token",
			mockAuthSetup: func() {
				mockAuthHelper.EXPECT().ValidateToken(invalidToken).Return(nil, errors.New("invalid token"))
			},
			mockRepoSetup:  func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockAuthSetup()
			tt.mockRepoSetup()

			req := &pb.ValidateRequest{Token: tt.token}
			resp, err := s.Validate(context.Background(), req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, int(resp.Status))
			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, resp.Error)
			}
		})
	}
}
