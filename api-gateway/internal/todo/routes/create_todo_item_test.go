package routes

import (
	"bytes"
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

func TestCreateTodoItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name                 string
		mockClient           *mocks.MockTodoServiceClient
		inputBody            string
		inputTodo            CreateTodoItemInput
		expectedStatusCode   int
		expectedResponseBody string
		userId               int64
	}{
		{
			name: "successfully creating todo item",
			mockClient: &mocks.MockTodoServiceClient{
				CreateTodoItemFunc: func(ctx context.Context, req *pb.CreateTodoItemRequest) (*pb.CreateTodoItemResponse, error) {
					return &pb.CreateTodoItemResponse{
						Status: http.StatusCreated,
					}, nil
				},
			},
			inputBody:            `{"title":"test todo item","description":"test description"}`,
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"status":201}`,
			userId:               1,
		},
		{
			name: "error creating todo item",
			mockClient: &mocks.MockTodoServiceClient{
				CreateTodoItemFunc: func(ctx context.Context, req *pb.CreateTodoItemRequest) (*pb.CreateTodoItemResponse, error) {
					return &pb.CreateTodoItemResponse{
						Status: http.StatusBadRequest,
						Error:  "invalid input body",
					}, nil
				},
			},
			inputBody:            "invalid body",
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid input body"}`,
			userId:               1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()

			req, _ := http.NewRequest(http.MethodPost, "/todo/1", bytes.NewBufferString(tt.inputBody))
			req.Header.Set("Content-Type", "application/json")

			r.POST("/todo/:id", func(ctx *gin.Context) {
				if tt.userId != 0 {
					ctx.Set(auth.Key, tt.userId)
				}
				CreateTodoItem(ctx, tt.mockClient)
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
