package main

import (
	"context"
	"log"

	"github.com/manarakozhamuratova/one-lab-task2/config"
	_ "github.com/manarakozhamuratova/one-lab-task2/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/manarakozhamuratova/one-lab-task2/internal/service"
	"github.com/manarakozhamuratova/one-lab-task2/internal/storage"
	"github.com/manarakozhamuratova/one-lab-task2/transport/httpserver"
	"github.com/manarakozhamuratova/one-lab-task2/transport/httpserver/handler"
	"github.com/manarakozhamuratova/one-lab-task2/transport/httpserver/middleware"
)

// @title Super API
// @version 1.0
// @description This is my first swagger documentation.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @host localhost:9090

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	cfg, err := config.ParseYAML()
	if err != nil {
		return err
	}
	if err := cfg.Validate(); err != nil {
		return err
	}
	st, err := storage.New(ctx, cfg)
	if err != nil {
		log.Fatal("storage init failed: ", err)
	}

	manager, err := service.NewManager(st)
	if err != nil {
		log.Fatal("manager init failed: ", err)
	}

	jwt := middleware.NewJWTAuth(cfg, manager.User)
	handlers := handler.NewHandler(cfg, manager, jwt)
	server := httpserver.NewServer(cfg, handlers)
	server.StartHTTPServer(ctx)
	return nil
}
