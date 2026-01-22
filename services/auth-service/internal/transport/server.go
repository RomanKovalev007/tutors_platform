package transport

import (
	"auth_service/pkg/kafka"
	"auth_service/internal/repository"
	"auth_service/internal/service"
	"auth_service/pkg/token"
	"database/sql"

	"github.com/redis/go-redis/v9"
	pb "github.com/RomanKovalev007/tutors_platform/api/gen/go/auth"
)


type ApiServer struct{
	pb.AuthServiceServer

	userService UserService
	authService AuthService
	
}

func NewApiServer(pgDB *sql.DB, redisClient *redis.Client, tokenCfg token.TokenConfig, producer *kafka.Producer) *ApiServer{
	userRepo := repository.NewUserRepository(pgDB)
	tokenRepo := repository.NewRedisTokenRepository(redisClient, "auth:refresh_token:")
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(tokenCfg, userRepo, tokenRepo, producer)
	
	return &ApiServer{
		userService: userService,
		authService: authService,
	}
}