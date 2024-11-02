package todo

import (
	"github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/config"
	pb "github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/todo/pb"
	"google.golang.org/grpc"
	"log"
)

type ServiceClient struct {
	Client pb.TodoServiceClient
}

func InitServiceClient(cfg *config.Config) pb.TodoServiceClient {
	cc, err := grpc.Dial(cfg.TodoSuvURL, grpc.WithInsecure())
	if err != nil {
		log.Println("Could not connect:", err)
	}

	log.Println("API Gateway - Todo Service Client initialized")
	return pb.NewTodoServiceClient(cc)
}
