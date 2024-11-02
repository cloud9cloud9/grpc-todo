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

type UpdateTodoListInput struct {
	Title string `json:"title"`
}

func UpdateTodoList(ctx *gin.Context, client pb.TodoServiceClient) {
	var req UpdateTodoListInput

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
		utils.NewErrorResponse(ctx, http.StatusBadRequest, invalidItemId)
		return
	}

	res, err := client.UpdateTodoList(context.Background(), &pb.UpdateTodoListRequest{
		UserId: userID,
		Id:     int64(listId),
		Title:  req.Title,
	})

	if err != nil {
		utils.NewErrorResponse(ctx, http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
