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

func DeleteTodoList(ctx *gin.Context, client pb.TodoServiceClient) {
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

	res, err := client.DeleteTodoList(context.Background(), &pb.DeleteTodoListRequest{
		Id:     int64(listId),
		UserId: userID,
	})

	if err != nil {
		utils.NewErrorResponse(ctx, http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
