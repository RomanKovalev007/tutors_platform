package main

import (
	"api_gateway/internal/clients/auth"
	"api_gateway/internal/config"
	api_http "api_gateway/internal/controller/http"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	cfg, err := config.ParseConfigFromEnv()
	if err != nil {
		panic(fmt.Errorf("could not parse config: %w", err))
	}

	authClient, err := auth.NewClient(cfg.AuthGRPC)
	if err != nil {
		panic(fmt.Errorf("failed to create auth client: %w", err))
	}
	defer authClient.Close()

	server := api_http.NewServer(authClient, cfg)

	server.RegisterHandlers()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.Start(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	graceSh := make(chan os.Signal, 1)
	signal.Notify(graceSh, os.Interrupt, syscall.SIGTERM)
	<-graceSh

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.GHTimeout)
	defer cancel()

	if err := server.Stop(shutdownCtx); err != nil {
		panic(err)
	}

	wg.Wait()
}
