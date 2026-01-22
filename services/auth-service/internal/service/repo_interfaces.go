package service

import (
	"auth_service/internal/models"
	"context"
	"time"
)

type UserRepository interface {
    CreateUser(ctx context.Context, user *models.User) (*models.User, error)
    DeleteUser(ctx context.Context, id string) error
    GetUserByID(ctx context.Context, id string) (*models.User, error)
    GetUserByEmail(ctx context.Context, email string) (*models.User, error)
    SetIsActive(ctx context.Context, id string, status bool) error
    UpdatePassword(ctx context.Context, id string, password string) error
	SelectAllUsers(ctx context.Context, limit, offset int32, isActive bool) ([]models.User, error)
}

type RefreshTokenRepository interface {
	DeleteToken(ctx context.Context, token string) error
	GetToken(ctx context.Context, token string) (string, error)
	SaveToken(ctx context.Context, userID string, token string, expiresAt time.Time) error
}