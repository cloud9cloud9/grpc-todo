package api

import (
	"github.com/cloud9cloud9/go-grpc-todo/auth-svc/internal/config"
	pb "github.com/cloud9cloud9/go-grpc-todo/auth-svc/internal/pb"
	"github.com/cloud9cloud9/go-grpc-todo/auth-svc/internal/service"
	"github.com/cloud9cloud9/go-grpc-todo/auth-svc/pkg/db"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	cfg *config.Config
}

func NewServer(
	cfg *config.Config,
) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) Start() error {
	repo := db.ConnectToPostgreSQL(s.cfg)
	log.Println("Database connected")

	lis, err := net.Listen("tcp", s.cfg.Port)
	if err != nil {
		log.Fatalln("failed at listening : ", err)
	}
	log.Println("Auth service started")
	serv := service.Server{
		Repo: repo,
	}
	log.Println("Server created")

	grpcServ := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServ, &serv)
	return grpcServ.Serve(lis)
}