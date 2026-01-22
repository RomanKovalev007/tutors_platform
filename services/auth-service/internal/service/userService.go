package service

import (
	"auth_service/internal/models"
	"auth_service/internal/repository"
	"context"
	"errors"
)

type userService struct {
    userRepo  UserRepository
}

func NewUserService(userRepo UserRepository) *userService {
    return &userService{
        userRepo:  userRepo,
    }
}

func (s *userService) GetAllUsers(ctx context.Context, limit, offset int32, isActive bool) ([]models.User, *models.Error){
	users, err := s.userRepo.SelectAllUsers(ctx, limit, offset, isActive)
	if err != nil{
		return nil, &models.Error{
			Code: models.INTERNALERROR,
			Message: err,
		}
	}

	return users, nil
} 

func (s *userService) GetUserByID(ctx context.Context, id string) (*models.User, *models.Error) {
    user, err := s.userRepo.GetUserByID(ctx, id)
    if err != nil {
        if errors.Is(err, repository.ErrUserNotFound){
			return nil, &models.Error{
				Code: models.USERNOTFOUND,
				Message: err,
			}
		}

		return nil, &models.Error{
			Code: models.INTERNALERROR,
			Message: err,
		}
    }

	return user, nil
}


func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, *models.Error) {
    user, err := s.userRepo.GetUserByEmail(ctx, email)
    if err != nil {
        if errors.Is(err, repository.ErrUserNotFound){
			return nil, &models.Error{
				Code: models.USERNOTFOUND,
				Message: err,
			}
		}

		return nil, &models.Error{
			Code: models.INTERNALERROR,
			Message: err,
		}
    }

	return user, nil
}

func (s *userService) UpdateIsActiveStatus(ctx context.Context, id string, status bool) *models.Error {
    err := s.userRepo.SetIsActive(ctx, id, status)
    if err != nil {
        if errors.Is(err, repository.ErrUserNotFound){
			return &models.Error{
				Code: models.USERNOTFOUND,
				Message: err,
			}
		}

		return &models.Error{
			Code: models.INTERNALERROR,
			Message: err,
		}
    }

	return nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) *models.Error {
    err := s.userRepo.DeleteUser(ctx, id)
    if err != nil {
        if errors.Is(err, repository.ErrUserNotFound){
			return &models.Error{
				Code: models.USERNOTFOUND,
				Message: err,
			}
		}

		return &models.Error{
			Code: models.INTERNALERROR,
			Message: err,
		}
    }

	return nil
}
