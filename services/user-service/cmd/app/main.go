package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"user-service/internal/config"
	kfk "user-service/pkg/kafka"
	"user-service/internal/transport"
	"user-service/pkg/migrator"
	"user-service/pkg/postgres"

	pb "github.com/RomanKovalev007/tutors_platform/api/gen/go/user"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)


func main() {
	cfg, err := config.ParseConfigFromEnv()
	if err != nil {
		panic(fmt.Errorf("could not parse config: %w", err))
	}
	cfg.DSN = cfg.FormatConnectionString()

	log.Println("config parse successful", cfg)

	err = migrator.RunMigrations(cfg.MigrationPath, cfg.DSN)
	if err != nil{
		panic(fmt.Errorf("failed to run migrations: %w", err))
	}

	pgDB, err := postgres.NewPostgresDB(cfg.PGConfig)
	if err != nil {
		panic(fmt.Errorf("failed to connect to postgres: %w", err))
	}
	defer pgDB.DB.Close()


	producer := kfk.NewProducer([]string{cfg.KafkaConfig.Brokers}, cfg.KafkaConfig.Topic)
	apiServer := transport.NewApiServer(pgDB.DB, producer)

	topics := []string{
		cfg.Topic,
		fmt.Sprintf("%s.retry.1", cfg.Topic),
		fmt.Sprintf("%s.retry.2", cfg.Topic),
		fmt.Sprintf("%s.retry.3", cfg.Topic),
	}

	ctx := context.Background()

	for _, topic := range topics {
		reader := kfk.NewConsumer([]string{cfg.KafkaConfig.Brokers}, topic, cfg.GroupID)

		go func(r *kafka.Reader) {
			for {
				msg, err := r.FetchMessage(ctx)
				if err != nil {
					continue
				}


				err = apiServer.EventHandler.Handle(ctx, msg)
				if err == nil {
					_ = r.CommitMessages(ctx, msg)
				}
			}
		}(reader)
	}


	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, apiServer)

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
