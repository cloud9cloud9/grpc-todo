package main

import (
	"github.com/cloud9cloud9/go-grpc-todo/auth-svc/cmd/api"
	"github.com/cloud9cloud9/go-grpc-todo/auth-svc/internal/config"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Could not load config : ", err)
	}

	server := api.NewServer(cfg)
	if err := server.Start(); err != nil {
		log.Fatalln("Could not start server : ", err)
	}
}
