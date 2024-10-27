package todo

import (
	"github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/auth"
	"github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/config"
	"github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/todo/routes"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, cfg *config.Config, authSvc *auth.ServiceClient) {
	a := auth.InitMiddleware(authSvc)

	svc := &ServiceClient{
		Client: InitServiceClient(cfg),
	}

	api := router.Group("/api", a.UserIdentity)
	{
		lists := api.Group("/lists")
		{
			lists.POST("/", svc.createTodoList)
			lists.GET("/", svc.getTodoLists)
			lists.GET("/:id", svc.getTodoListById)
			lists.PUT("/:id", svc.updateTodoList)
			lists.DELETE("/:id", svc.deleteTodoList)

			items := lists.Group(":id/items")
			{
				items.POST("/", svc.createTodoItem)
				items.GET("/", svc.getTodoItems)
			}
		}

		items := api.Group("items")
		{
			items.GET("/:id", svc.getTodoItemById)
			items.PUT("/:id", svc.updateTodoItem)
			items.DELETE("/:id", svc.deleteTodoItemById)
		}
	}
}

func (svc *ServiceClient) createTodoList(ctx *gin.Context) {
	routes.CreateTodoList(ctx, svc.Client)
}

func (svc *ServiceClient) getTodoLists(ctx *gin.Context) {
	routes.GetTodoLists(ctx, svc.Client)
}

func (svc *ServiceClient) getTodoListById(ctx *gin.Context) {
	routes.GetTodoListById(ctx, svc.Client)
}

func (svc *ServiceClient) updateTodoList(ctx *gin.Context) {
	routes.UpdateTodoList(ctx, svc.Client)
}

func (svc *ServiceClient) deleteTodoList(ctx *gin.Context) {
	routes.DeleteTodoList(ctx, svc.Client)
}

func (svc *ServiceClient) createTodoItem(ctx *gin.Context) {
	routes.CreateTodoItem(ctx, svc.Client)
}

func (svc *ServiceClient) getTodoItems(ctx *gin.Context) {
	routes.GetTodoItems(ctx, svc.Client)
}

func (svc *ServiceClient) getTodoItemById(ctx *gin.Context) {
	routes.GetTodoItemById(ctx, svc.Client)
}

func (svc *ServiceClient) updateTodoItem(ctx *gin.Context) {
	routes.UpdateTodoItem(ctx, svc.Client)
}

func (svc *ServiceClient) deleteTodoItemById(ctx *gin.Context) {
	routes.DeleteTodoItemById(ctx, svc.Client)
}
