package routes

import (
	"context"
	"github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/auth"
	pb "github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/todo/pb"
	"github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/todo/routes/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetTodoItems(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name                 string
		mockClient           *mocks.MockTodoServiceClient
		expectedStatusCode   int
		expectedResponseBody string
		userId               int64
	}{
		{
			name: "successfully getting todo items",
			mockClient: &mocks.MockTodoServiceClient{
				GetTodoItemsFunc: func(ctx context.Context, req *pb.GetTodoItemsRequest) (*pb.GetTodoItemsResponse, error) {
					return &pb.GetTodoItemsResponse{
						Status: http.StatusOK,
					}, nil
				},
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":200}`,
			userId:               1,
		},
		{
			name: "error getting todo items",
			mockClient: &mocks.MockTodoServiceClient{
				GetTodoItemsFunc: func(ctx context.Context, req *pb.GetTodoItemsRequest) (*pb.GetTodoItemsResponse, error) {
					return &pb.GetTodoItemsResponse{
						Status: http.StatusBadGateway,
						Error:  "invalid input body",
					}, nil
				},
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":502,"error":"invalid input body"}`,
			userId:               1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()

			req, _ := http.NewRequest(http.MethodGet, "/items/1", nil)
			req.Header.Set("Content-Type", "application/json")

			r.GET("/items/:id", func(ctx *gin.Context) {
				if tt.userId != 0 {
					ctx.Set(auth.Key, tt.userId)
				}
				GetTodoItems(ctx, tt.mockClient)
			})

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)

			if tt.expectedResponseBody != "" {
				assert.Equal(t, tt.expectedResponseBody, strings.TrimSpace(w.Body.String()))
			}
		})
	}
}
