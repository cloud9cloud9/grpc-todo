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

func TestCreateTodoList(t *testing.T) {
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
			name: "successfully creating todo list",
			mockClient: &mocks.MockTodoServiceClient{
				CreateTodoListFunc: func(ctx context.Context, req *pb.CreateTodoListRequest) (*pb.CreateTodoListResponse, error) {
					return &pb.CreateTodoListResponse{
						Status: http.StatusCreated,
					}, nil
				},
			},
			inputBody:            `{"title":"test todo list"}`,
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"status":201}`,
			userId:               1,
		},
		{
			name: "error creating todo list",
			mockClient: &mocks.MockTodoServiceClient{
				CreateTodoListFunc: func(ctx context.Context, req *pb.CreateTodoListRequest) (*pb.CreateTodoListResponse, error) {
					return &pb.CreateTodoListResponse{
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

			req, _ := http.NewRequest(http.MethodPost, "/list", bytes.NewBufferString(tt.inputBody))
			req.Header.Set("Content-Type", "application/json")

			r.POST("/list", func(ctx *gin.Context) {
				if tt.userId != 0 {
					ctx.Set(auth.Key, tt.userId)
				}
				CreateTodoList(ctx, tt.mockClient)
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
