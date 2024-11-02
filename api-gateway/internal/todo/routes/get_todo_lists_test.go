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
	"testing"
)

func TestGetTodoLists(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name                 string
		mockClient           *mocks.MockTodoServiceClient
		expectedStatusCode   int
		expectedResponseBody string
		userId               int64
	}{
		{
			name: "successfully getting todo lists",
			mockClient: &mocks.MockTodoServiceClient{
				GetTodoListsFunc: func(ctx context.Context, req *pb.GetTodoListsRequest) (*pb.GetTodoListsResponse, error) {
					return &pb.GetTodoListsResponse{
						Status: http.StatusOK,
					}, nil
				},
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":200}`,
			userId:               1,
		},
		{
			name: "error getting todo lists",
			mockClient: &mocks.MockTodoServiceClient{
				GetTodoListsFunc: func(ctx context.Context, req *pb.GetTodoListsRequest) (*pb.GetTodoListsResponse, error) {
					return &pb.GetTodoListsResponse{
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

			req, _ := http.NewRequest(http.MethodGet, "/lists", nil)
			req.Header.Set("Content-Type", "application/json")

			r.GET("/lists", func(ctx *gin.Context) {
				if tt.userId != 0 {
					ctx.Set(auth.Key, tt.userId)
				}

				GetTodoLists(ctx, tt.mockClient)
			})

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}
