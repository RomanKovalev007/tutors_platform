package transport

import (
	"auth_service/internal/models"
	"context"
)

type UserService interface{
	DeleteUser(ctx context.Context, id string) *models.Error
	GetAllUsers(ctx context.Context, limit int32, offset int32, isActive bool) ([]models.User, *models.Error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, *models.Error)
	GetUserByID(ctx context.Context, id string) (*models.User, *models.Error)
	UpdateIsActiveStatus(ctx context.Context, id string, status bool) *models.Error
}

type AuthService interface{
	Register(ctx context.Context, email string, password string, confirmPassword string) (*models.User, *models.Tokens, *models.Error)
	Login(ctx context.Context, email string, password string) (*models.User, *models.Tokens, *models.Error)
	Logout(ctx context.Context, token string) *models.Error
	Refresh(ctx context.Context, rt string) (*models.Tokens, *models.Error)
	ValidateAccessToken(ctx context.Context, tokenStr string) (string, int64, *models.Error)
}
