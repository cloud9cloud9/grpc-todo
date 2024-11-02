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

const (
	testEmail     = "test@example.com"
	testPassword  = "password"
	testToken     = "fake-token"
	wrongEmail    = "wrong-email"
	wrongPassword = "wrong-password"
)

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockClient     *mocks.MockAuthServiceClient
		reqBody        LoginRequest
		expectedStatus int
		expectedToken  string
	}{
		{
			name: "valid credentials",
			mockClient: &mocks.MockAuthServiceClient{
				LoginFunc: func(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
					return &pb.LoginResponse{Token: testToken}, nil
				},
			},
			reqBody:        LoginRequest{Email: testEmail, Password: testPassword},
			expectedStatus: http.StatusCreated,
			expectedToken:  testToken,
		},
		{
			name: "invalid credentials",
			mockClient: &mocks.MockAuthServiceClient{
				LoginFunc: func(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
					return nil, status.Error(codes.Unauthenticated, "invalid credentials")
				},
			},
			reqBody:        LoginRequest{Email: wrongEmail, Password: wrongPassword},
			expectedStatus: http.StatusBadGateway,
			expectedToken:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/login", func(ctx *gin.Context) {
				Login(ctx, tt.mockClient)
			})

			body, _ := json.Marshal(tt.reqBody)
			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus == http.StatusCreated {
				var res pb.LoginResponse
				err := json.Unmarshal(rec.Body.Bytes(), &res)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedToken, res.Token)
			} else {
				assert.Equal(t, tt.expectedStatus, http.StatusBadGateway)
				assert.Equal(t, tt.expectedToken, "")
			}
		})
	}
}
