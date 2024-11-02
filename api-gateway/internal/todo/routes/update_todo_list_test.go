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

func TestUpdateTodoList(t *testing.T) {
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
			name: "successfully updating todo list",
			mockClient: &mocks.MockTodoServiceClient{
				UpdateTodoListFunc: func(ctx context.Context, req *pb.UpdateTodoListRequest) (*pb.UpdateTodoListResponse, error) {
					return &pb.UpdateTodoListResponse{
						Status: http.StatusOK,
						List: &pb.TodoList{
							Id:    1,
							Title: "title",
						},
					}, nil
				},
			},
			inputBody:            `{"id":"1","userId":"1","title":"title"}`,
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"list":{"id":1,"title":"title"},"status":200}`,
			userId:               1,
		},
		{
			name: "error updating todo list",
			mockClient: &mocks.MockTodoServiceClient{
				UpdateTodoListFunc: func(ctx context.Context, req *pb.UpdateTodoListRequest) (*pb.UpdateTodoListResponse, error) {
					return &pb.UpdateTodoListResponse{
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

			req, _ := http.NewRequest(http.MethodPut, "/list/1", strings.NewReader(tt.inputBody))
			req.Header.Set("Content-Type", "application/json")

			r.PUT("/list/:id", func(ctx *gin.Context) {
				if tt.userId != 0 {
					ctx.Set(auth.Key, tt.userId)
				}

				UpdateTodoList(ctx, tt.mockClient)
			})

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}
