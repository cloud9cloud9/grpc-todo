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

func TestDeleteTodoList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name                 string
		mockClient           *mocks.MockTodoServiceClient
		inputBody            string
		expectedStatusCode   int
		expectedResponseBody string
		userId               int64
	}{
		{
			name: "successfully deleting todo list",
			mockClient: &mocks.MockTodoServiceClient{
				DeleteTodoListFunc: func(ctx context.Context, req *pb.DeleteTodoListRequest) (*pb.DeleteTodoListResponse, error) {
					return &pb.DeleteTodoListResponse{
						Status: http.StatusOK,
					}, nil
				},
			},
			inputBody:            `{"id":"1"","userId":"1""}`,
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":200}`,
			userId:               1,
		},
		{
			name: "error deleting todo list",
			mockClient: &mocks.MockTodoServiceClient{
				DeleteTodoListFunc: func(ctx context.Context, req *pb.DeleteTodoListRequest) (*pb.DeleteTodoListResponse, error) {
					return &pb.DeleteTodoListResponse{
						Status: http.StatusBadGateway,
						Error:  "invalid input body",
					}, nil
				},
			},
			inputBody:            `{"userId":1}`,
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":502,"error":"invalid input body"}`,
			userId:               1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()

			req, _ := http.NewRequest(http.MethodDelete, "/todo-list/1", bytes.NewBufferString(tt.inputBody))
			req.Header.Set("Content-Type", "application/json")

			r.DELETE("/todo-list/:id", func(ctx *gin.Context) {
				if tt.userId != 0 {
					ctx.Set(auth.Key, tt.userId)
				}
				DeleteTodoList(ctx, tt.mockClient)
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
