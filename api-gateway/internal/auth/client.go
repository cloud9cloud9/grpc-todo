package auth

import (
	pb "github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/auth/pb"
	"github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/config"
	"google.golang.org/grpc"
	"log"
)

type ServiceClient struct {
	Client pb.AuthServiceClient
}

func InitServiceClient(cfg *config.Config) pb.AuthServiceClient {
	cc, err := grpc.Dial(cfg.AuthSuvURL, grpc.WithInsecure())
	if err != nil {
		log.Println("Could not connect:", err)
	}

	log.Println("API Gateway - Auth Service Client initialized")
	return pb.NewAuthServiceClient(cc)
}
