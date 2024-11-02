package routes

import (
	"context"
	"github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/auth"
	pb "github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/todo/pb"
	"github.com/cloud9cloud9/go-grpc-todo/api-gateway/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var (
	invalidListID = "invalid list id"
)

type CreateTodoItemInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func CreateTodoItem(ctx *gin.Context, client pb.TodoServiceClient) {
	var req CreateTodoItemInput

	if err := ctx.BindJSON(&req); err != nil {
		utils.NewErrorResponse(ctx, http.StatusBadRequest, invalidInputBody)
		return
	}

	userID, err := auth.GetUserId(ctx)
	if err != nil {
		utils.NewErrorResponse(ctx, http.StatusUnauthorized, invalidUserID)
		return
	}

	listId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.NewErrorResponse(ctx, http.StatusBadRequest, invalidListID)
		return
	}

	res, err := client.CreateTodoItem(context.Background(), &pb.CreateTodoItemRequest{
		ListId:      int64(listId),
		UserId:      userID,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		utils.NewErrorResponse(ctx, http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}
