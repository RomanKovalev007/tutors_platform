package service

import (
	"auth_service/internal/events"
	"auth_service/internal/models"
	"auth_service/internal/repository"
	"auth_service/pkg/token"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type EventPublisher interface {
	Publish(ctx context.Context, topic string, key string, value any) error
}

type authService struct {
    producer EventPublisher
    topic    string
    userRepo  UserRepository
    tokenRepo  RefreshTokenRepository
    tokens *token.Manager
}

func NewAuthService(tokenCfg token.TokenConfig, userRepo UserRepository, tokenRepo RefreshTokenRepository, producer EventPublisher, topic string) *authService {
    if producer == nil {
		panic("kafka producer is nil")
	}
    return &authService{
        userRepo:  userRepo,
        tokenRepo:  tokenRepo,
        tokens: token.NewManager(tokenCfg),
        producer: producer,
        topic:    topic,
    }
}

func (s *authService) Register(ctx context.Context, email, password, confirmPassword string) (*models.User, *models.Tokens, *models.Error) {
    if  password != confirmPassword{
        return nil, nil, &models.Error{Code: models.INVALIDINPUT, Message: fmt.Errorf("password not equal")}
    }
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, nil, &models.Error{Code: models.INVALIDINPUT, Message: err}
    }

    req_user := &models.User{
        ID:       generateID(),
        Email:    email,
        Password: string(hash),
    }

    user, err := s.userRepo.CreateUser(ctx, req_user)
    if err != nil{
        if errors.Is(err, repository.ErrUserExists){
            return nil, nil, &models.Error{Code: models.USEREXISTS, Message: err}
        }
        return nil, nil, &models.Error{Code: models.INTERNALERROR, Message: err}
    }

    tokens, err := s.newTokens(ctx, user.ID)
    if err != nil{
        return nil, nil, &models.Error{Code: models.INTERNALERROR, Message: err}
    }

    event := events.UserRegisteredEvent{
        EventType:  "UserRegistered",
        EventID:    uuid.New().String(),
        OccurredAt: time.Now().UTC(),
        Payload: events.UserRegisteredPayload{
            UserID: user.ID,
            Email:  email,
    },
	}


	if err = s.producer.Publish(ctx, s.topic, user.ID, event); err != nil{
        log.Println("cant publish user registred event")
    }

    
    return user, tokens, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*models.User, *models.Tokens, *models.Error) {
    user, err := s.userRepo.GetUserByEmail(ctx, email)
    if err != nil {
        if errors.Is(err, repository.ErrUserNotFound){
            return nil, nil, &models.Error{Code: models.INVALIDCREDENTIALS, Message: err}
        }
        return nil, nil, &models.Error{Code: models.INTERNALERROR, Message: err}
    }

    if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
        return nil, nil, &models.Error{Code: models.INVALIDCREDENTIALS, Message: fmt.Errorf("invalid password")}
    }

    tokens, err := s.newTokens(ctx, user.ID)
    if err != nil{
        return nil, nil, &models.Error{Code: models.INTERNALERROR, Message: err}
    }

    return user, tokens, nil
}

func (s *authService) Logout(ctx context.Context, token string) *models.Error {
    if err := s.tokenRepo.DeleteToken(ctx, token); err != nil{
        if errors.Is(err, repository.ErrUserNotFound){
            return &models.Error{Code: models.INVALIDTOKEN, Message: err}
        }
        return &models.Error{Code: models.INTERNALERROR, Message: err}
    }
    return nil
}



func (s *authService) Refresh(ctx context.Context, rt string) (*models.Tokens, *models.Error) {
    userID, err := s.tokenRepo.GetToken(ctx, rt)
    if err != nil {
        if errors.Is(err, repository.ErrUserNotFound){
            return nil, &models.Error{Code: models.INVALIDTOKEN, Message: err}
        }
        return nil, &models.Error{Code: models.INTERNALERROR, Message: err}
    }

    err = s.tokenRepo.DeleteToken(ctx, rt)
    if err != nil{
        return nil, &models.Error{Code: models.INTERNALERROR, Message: err}
    }


    tokens, err := s.newTokens(ctx, userID)
    if err != nil{
        return nil, &models.Error{Code: models.INTERNALERROR, Message: err}
    }

    return tokens, nil
}

func (s *authService) ValidateAccessToken(ctx context.Context, tokenStr string) (string, int64, *models.Error) {
    userID, exp, err := s.tokens.ValidateAccessToken(tokenStr)
    if err != nil {
        if errors.Is(err, token.ErrTokenExpired){
            return userID, 0, &models.Error{Code: models.TOKENEXPIRED, Message: err}
        }

        return "", 0, &models.Error{Code: models.INVALIDTOKEN, Message: err}
    }

    return userID, exp, nil
}

func (s *authService) newTokens(ctx context.Context, userID string) (*models.Tokens, error) {
    access, err := s.tokens.NewAccessToken(userID)
    if err != nil {
        return nil, err
    }

    refresh, exp, err := s.tokens.NewRefreshToken()
    if err != nil {
        return nil, err
    }

    if err := s.tokenRepo.SaveToken(ctx, userID, refresh, exp); err != nil {
        return nil, err
    }

    return &models.Tokens{
        AccessToken:  access,
        RefreshToken: refresh,
        ExpiresAt:    int32(exp.Unix()),
    }, nil
}



