package main

import (
	"auth_service/internal/config"
	"auth_service/internal/transport"
	"auth_service/pkg/kafka"
	"auth_service/pkg/migrator"
	"auth_service/pkg/postgres"
	"auth_service/pkg/redis"
	"auth_service/pkg/token"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	pb "github.com/RomanKovalev007/tutors_platform/api/gen/go/auth"

	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.ParseConfigFromEnv()
	if err != nil {
		panic(fmt.Errorf("could not parse config: %w", err))
	}
	cfg.DSN = cfg.FormatConnectionString()

	log.Println("config parse successful", cfg)

	producer := kafka.NewProducer([]string{cfg.KafkaConfig.Brokers}, cfg.KafkaConfig.Topic)
	if producer == nil {
		panic("failed to create kafka producer")
	}
	defer producer.Close()

	err = migrator.RunMigrations(cfg.MigrationPath, cfg.DSN)
	if err != nil {
		panic(fmt.Errorf("failed to run migrations: %w", err))
	}

	pgDB, err := postgres.NewPostgresDB(cfg.PGConfig)
	if err != nil {
		panic(fmt.Errorf("failed to connect to postgres: %w", err))
	}
	defer pgDB.DB.Close()

	redisDB, err := redis.NewRedisDB(cfg.RedisConfig)
	if err != nil {
		panic(fmt.Errorf("failed to connect to redis: %w", err))
	}
	defer redisDB.Close()

	tokenCfg := token.TokenConfig{}
	tokenCfg.AccessTTL = time.Duration(cfg.AccessTTL) * time.Minute
	tokenCfg.RefreshTTL = time.Duration(cfg.RefreshTTL) * time.Hour

	apiServer := transport.NewApiServer(pgDB.DB, redisDB, tokenCfg, producer, cfg.KafkaConfig.Topic)

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, apiServer)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		panic(fmt.Errorf("error starting tcp listener: %w", err))
	}

	log.Println("tcp listener started")

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			panic(fmt.Errorf("error grpc server serving: %w", err))
		}
	}()

	log.Println("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	grpcServer.GracefulStop()
	log.Println("server shut down")

}
