package api

import (
	"github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/auth"
	"github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/config"
	"github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/todo"
	"github.com/gin-gonic/gin"
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
	router := gin.Default()
	authRoutes := auth.RegisterRoutes(router, s.cfg)
	todo.RegisterRoutes(router, s.cfg, authRoutes)
	return router.Run(s.cfg.Port)
}
