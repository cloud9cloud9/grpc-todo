package api

import (
	"github.com/cloud9cloud9/go-grpc-todo/todo-svc/internal/config"
	pb "github.com/cloud9cloud9/go-grpc-todo/todo-svc/internal/pb"
	"github.com/cloud9cloud9/go-grpc-todo/todo-svc/internal/repository"
	"github.com/cloud9cloud9/go-grpc-todo/todo-svc/internal/service"
	"github.com/cloud9cloud9/go-grpc-todo/todo-svc/pkg/db"
	"github.com/cloud9cloud9/go-grpc-todo/todo-svc/pkg/utils"
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
	database := db.ConnectToPostgreSQL(s.cfg)
	repo := repository.NewRepository(database)
	mapper := utils.NewMapper()
	log.Println("Database connected")

	lis, err := net.Listen("tcp", s.cfg.Server.Port)
	if err != nil {
		log.Fatalln("failed at listening : ", err)
	}
	log.Println("Auth service started")
	serv := service.Server{
		ListRepo: repo.TodoList,
		ItemRepo: repo.TodoItem,
		Mapper:   mapper,
	}
	log.Println("Server created")

	grpcServ := grpc.NewServer()
	pb.RegisterTodoServiceServer(grpcServ, &serv)
	return grpcServ.Serve(lis)
}
