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

func TestUpdateTodoItem(t *testing.T) {
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
			name: "successfully updating todo item",
			mockClient: &mocks.MockTodoServiceClient{
				UpdateTodoItemFunc: func(ctx context.Context, req *pb.UpdateTodoItemRequest) (*pb.UpdateTodoItemResponse, error) {
					return &pb.UpdateTodoItemResponse{
						Status: http.StatusOK,
						Item: &pb.TodoItem{
							Id:          1,
							Title:       "title",
							Description: "description",
							Completed:   true,
						},
					}, nil
				},
			},
			inputBody:            `{"id":"1","userId":"1","title":"title","description":"description","done":true}`,
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"item":{"id":1,"title":"title","description":"description","completed":true},"status":200}`,
			userId:               1,
		},
		{
			name: "error updating todo item",
			mockClient: &mocks.MockTodoServiceClient{
				UpdateTodoItemFunc: func(ctx context.Context, req *pb.UpdateTodoItemRequest) (*pb.UpdateTodoItemResponse, error) {
					return &pb.UpdateTodoItemResponse{
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

			req, _ := http.NewRequest(http.MethodPut, "/update-item/1", strings.NewReader(tt.inputBody))
			req.Header.Set("Content-Type", "application/json")

			r.PUT("/update-item/:id", func(ctx *gin.Context) {
				if tt.userId != 0 {
					ctx.Set(auth.Key, tt.userId)
				}

				UpdateTodoItem(ctx, tt.mockClient)
			})

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}
