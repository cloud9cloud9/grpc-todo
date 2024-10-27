package routes

import (
	"context"
	pb "github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/auth/pb"
	"github.com/cloud9cloud9/go-grpc-todo/api-gateway/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	invalidInputBody = "invalid input body"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(ctx *gin.Context, client pb.AuthServiceClient) {
	var req LoginRequest

	if err := ctx.BindJSON(&req); err != nil {
		utils.NewErrorResponse(ctx, http.StatusBadRequest, invalidInputBody)
		return
	}

	res, err := client.Login(context.Background(), &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		utils.NewErrorResponse(ctx, http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}
