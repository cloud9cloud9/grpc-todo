package routes

import (
	"github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/auth"
	pb "github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/todo/pb"
	"github.com/cloud9cloud9/go-grpc-todo/api-gateway/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var (
	invalidItemId = "invalid must be integer"
)

func DeleteTodoItemById(ctx *gin.Context, c pb.TodoServiceClient) {
	itemId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.NewErrorResponse(ctx, http.StatusBadRequest, invalidItemId)
		return
	}

	userID, err := auth.GetUserId(ctx)
	if err != nil {
		utils.NewErrorResponse(ctx, http.StatusUnauthorized, invalidUserID)
		return
	}

	res, err := c.DeleteTodoItem(ctx, &pb.DeleteTodoItemRequest{
		UserId: userID,
		Id:     int64(itemId),
	})

	if err != nil {
		utils.NewErrorResponse(ctx, http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
