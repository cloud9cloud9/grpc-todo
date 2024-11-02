package auth

import (
	"context"
	"encoding/json"
	pb "github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/auth/pb"
	"github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/auth/routes/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	validToken   = "Bearer valid-token"
	invalidToken = "invalid-token"
	missingToken = ""
)

func TestUserIdentity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		header         string
		mockClient     *mocks.MockAuthServiceClient
		expectedStatus int
		expectedUserId interface{}
	}{
		{
			name:   "valid token",
			header: validToken,
			mockClient: &mocks.MockAuthServiceClient{
				ValidateFunc: func(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
					return &pb.ValidateResponse{
						UserId: 1,
						Status: http.StatusOK,
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedUserId: 1,
		},
		{
			name:   "missing token",
			header: missingToken,
			mockClient: &mocks.MockAuthServiceClient{
				ValidateFunc: func(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
					return &pb.ValidateResponse{
						UserId: 0,
						Status: http.StatusUnauthorized,
					}, nil
				},
			},
			expectedStatus: http.StatusUnauthorized,
			expectedUserId: 0,
		},
		{
			name:   "invalid token format",
			header: invalidToken,
			mockClient: &mocks.MockAuthServiceClient{
				ValidateFunc: func(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
					return &pb.ValidateResponse{
						UserId: 0,
						Status: http.StatusUnauthorized,
					}, nil
				},
			},
			expectedStatus: http.StatusUnauthorized,
			expectedUserId: 0,
		},
		{
			name:   "invalid token",
			header: "Bearer invalid-token",
			mockClient: &mocks.MockAuthServiceClient{
				ValidateFunc: func(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
					return &pb.ValidateResponse{
						UserId: 0,
						Status: http.StatusUnauthorized,
					}, nil
				},
			},
			expectedStatus: http.StatusUnauthorized,
			expectedUserId: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			var serviceClient *ServiceClient
			if tt.mockClient != nil {
				serviceClient = &ServiceClient{Client: tt.mockClient}
			}
			middleware := InitMiddleware(serviceClient)
			router.Use(middleware.UserIdentity)
			router.GET("/protected", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{"user_id": ctx.Value(Key)})
			})

			req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
			if tt.header != "" {
				req.Header.Set(authHeader, tt.header)
			}
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUserId, 1)
			}
		})
	}
}
