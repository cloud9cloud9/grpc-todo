package routes

import (
	"bytes"
	"context"
	"encoding/json"
	pb "github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/auth/pb"
	"github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/auth/routes/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockClient     *mocks.MockAuthServiceClient
		reqBody        RegisterRequest
		expectedStatus int
		expectedError  string
	}{
		{
			name: "valid credentials",
			mockClient: &mocks.MockAuthServiceClient{
				RegisterFunc: func(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
					return &pb.RegisterResponse{Status: http.StatusCreated}, nil
				},
			},
			reqBody:        RegisterRequest{Email: testEmail, Password: testPassword},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			name: "invalid credentials",
			mockClient: &mocks.MockAuthServiceClient{
				RegisterFunc: func(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
					return nil, status.Error(codes.Unauthenticated, "invalid credentials")
				},
			},
			reqBody:        RegisterRequest{Email: wrongEmail, Password: wrongPassword},
			expectedStatus: http.StatusBadGateway,
			expectedError:  "rpc error: code = Unauthenticated desc = invalid credentials",
		},
		{
			name: "email already exists",
			mockClient: &mocks.MockAuthServiceClient{
				RegisterFunc: func(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
					return nil, status.Error(codes.AlreadyExists, "email already exists")
				},
			},
			reqBody:        RegisterRequest{Email: testEmail, Password: testPassword},
			expectedStatus: http.StatusBadGateway,
			expectedError:  "rpc error: code = AlreadyExists desc = email already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/register", func(ctx *gin.Context) {
				Register(ctx, tt.mockClient)
			})

			body, _ := json.Marshal(tt.reqBody)
			req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusBadGateway {
				var errorResponse map[string]string
				err := json.Unmarshal(rr.Body.Bytes(), &errorResponse)
				assert.NoError(t, err)

				assert.Contains(t, errorResponse["message"], tt.expectedError)
			}
		})
	}
}
