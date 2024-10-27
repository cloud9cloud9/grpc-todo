package auth

import (
	"github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/auth/routes"
	"github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/config"
	"github.com/gin-gonic/gin"
	"log"
)

func RegisterRoutes(router *gin.Engine, cfg *config.Config) *ServiceClient {
	log.Println("API Gateway - Auth Service Routes initialized")
	svc := &ServiceClient{
		Client: InitServiceClient(cfg),
	}

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", svc.Register)
		auth.POST("/sign-in", svc.Login)
	}

	return svc
}

func (svc *ServiceClient) Register(ctx *gin.Context) {
	routes.Register(ctx, svc.Client)
}
func (svc *ServiceClient) Login(ctx *gin.Context) {
	routes.Login(ctx, svc.Client)
}
