package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"task_service/internal/config"
	"task_service/internal/repository"
	"task_service/internal/service"
	"task_service/internal/transport"
	client "task_service/pkg/groupClient"
	"task_service/pkg/pool"
)

type app struct {
	server      *transport.Server
	repo        *repository.Repository
	groupClient *client.GroupClient
}

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	app := &app{}

	pool, err := pool.CreateNewPool(cfg.PostgresCfg)
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}

	repo := repository.NewRepository(pool)

	client, err := client.NewGroupClient(cfg.GroupServiceAddr)
	if err != nil {
		log.Printf("failed to create group client: %v", err) // заменить на фатал
	}

	service := service.NewService(repo, client)

	log.Printf("starting server on port %s...", cfg.ServerPort)
	server, err := transport.NewServer(cfg.ServerPort, service)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	app.groupClient = client
	app.repo = repo
	app.server = server

	app.server.Start()

	app.WaitForGracefulShutdown()
}

func (a *app) WaitForGracefulShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Printf("shutting down server...")
	a.server.Stop()

	log.Printf("closing group_service connection...")
	err := a.groupClient.Close()
	if err != nil {
		log.Printf("failed to close group_service connection: %v", err)
	}

	log.Printf("closing database connection...")
	a.repo.Close()

	log.Printf("app shutted down successfully")
}
