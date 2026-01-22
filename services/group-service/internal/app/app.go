package app

import (
	"context"
	"fmt"
	"group_service/internal/adapter"
	"group_service/internal/clients/user"
	"group_service/internal/config"
	"group_service/internal/controller/grpc"
	"group_service/internal/usecase"
	postgres "group_service/pkg/db"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type App struct {
	grpcServer *grpc.Server
	postgresDB *postgres.Database
}

func NewApp(cfg *config.Config) (*App, error) {
	db, err := postgres.New(cfg.PostgresConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	groupsRepo := adapter.NewGroupsRepo(db.Pool)

	userClient, err := user.NewClient(cfg.UserServiceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	groupsUsecase := usecase.NewGroupsUsecase(groupsRepo, userClient)

	server := grpc.NewServer(groupsUsecase)

	return &App{
		grpcServer: server,
		postgresDB: db,
	}, nil
}

func (a *App) MustRun(ctx context.Context, grpcPort int, timeout time.Duration) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.grpcServer.Run(grpcPort); err != nil {
			panic(err)
		}
	}()

	graceSh := make(chan os.Signal, 1)
	signal.Notify(graceSh, os.Interrupt, syscall.SIGTERM)
	<-graceSh

	a.grpcServer.Stop()

	a.postgresDB.Close()

	wg.Wait()
}
