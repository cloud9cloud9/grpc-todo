package routes

import (
	"context"
	"github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/auth"
	pb "github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/todo/pb"
	"github.com/cloud9cloud9/go-grpc-todo/api-gateway/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	invalidInputBody = "invalid input body"
	invalidUserID    = "invalid user ID"
)

type CreateTodoListInput struct {
	Title string `json:"title"`
}

func CreateTodoList(ctx *gin.Context, client pb.TodoServiceClient) {
	var req CreateTodoListInput

	if err := ctx.BindJSON(&req); err != nil {
		utils.NewErrorResponse(ctx, http.StatusBadRequest, invalidInputBody)
		return
	}

	userID, err := auth.GetUserId(ctx)
	if err != nil {
		utils.NewErrorResponse(ctx, http.StatusUnauthorized, invalidUserID)
		return
	}

	res, err := client.CreateTodoList(context.Background(), &pb.CreateTodoListRequest{
		UserId: userID,
		Title:  req.Title,
	})

	if err != nil {
		utils.NewErrorResponse(ctx, http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}
