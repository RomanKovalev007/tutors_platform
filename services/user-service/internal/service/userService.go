package service

import (
	"context"
	"errors"
	"user-service/internal/models"
	"user-service/internal/repository"
)

type UserService struct {
	userRepo    UserProfileRepository
	tutorRepo   TutorProfileRepository
	studentRepo StudentProfileRepository
}

func NewUserService(
	userRepo UserProfileRepository,
	tutorRepo TutorProfileRepository,
	studentRepo StudentProfileRepository,
) *UserService {
	return &UserService{
		userRepo:    userRepo,
		tutorRepo:   tutorRepo,
		studentRepo: studentRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.UserProfile) (*models.UserProfile, *models.Error) {
	
	resp_user, err := s.userRepo.CreateUser(ctx, user)
	if err != nil{
		if errors.Is(err, repository.ErrUserExists){
			return nil, &models.Error{Code: models.USEREXISTS, Message: err}
		}
		return nil, &models.Error{Code: models.INTERNALERROR, Message: err}
	}
	
	return resp_user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.UserProfile, *models.Error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil{
		if errors.Is(err, repository.ErrUserNotFound){
			return nil, &models.Error{Code: models.USERNOTFOUND, Message: err}
		}
		return nil, &models.Error{Code: models.INTERNALERROR, Message: err}
	}

	return user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.UserProfile, *models.Error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil{
		if errors.Is(err, repository.ErrUserNotFound){
			return nil, &models.Error{Code: models.USERNOTFOUND, Message: err}
		}
		return nil, &models.Error{Code: models.INTERNALERROR, Message: err}
	}

	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *models.UserProfile) (*models.UserProfile, *models.Error) {
	resp_user, err := s.userRepo.UpdateUser(ctx, user)
	if err != nil{
		if errors.Is(err, repository.ErrUserNotFound){
			return nil, &models.Error{Code: models.USERNOTFOUND, Message: err}
		}
		return nil, &models.Error{Code: models.INTERNALERROR, Message: err}
	}

	return resp_user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) *models.Error {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil{
		if errors.Is(err, repository.ErrUserNotFound){
			return &models.Error{Code: models.USERNOTFOUND, Message: err}
		}
		return &models.Error{Code: models.INTERNALERROR, Message: err}
	}

	if user.IsTutor {
		err = s.tutorRepo.DeleteTutorProfie(ctx, id)
		if err != nil{
			return &models.Error{Code: models.INTERNALERROR, Message: err}
		}
	}

	if user.IsStudent {
		err = s.studentRepo.DeleteStudentProfie(ctx, id)
		if err != nil{
			return &models.Error{Code: models.INTERNALERROR, Message: err}
		}
	}

	err = s.userRepo.DeleteUser(ctx, id)
	if err != nil{
		return &models.Error{Code: models.INTERNALERROR, Message: err}
	}

	return nil

}

func (s *UserService) GetUserTypes(ctx context.Context, id string) (*models.UserType, *models.Error) {
	utypes, err := s.userRepo.GetUserTypes(ctx, id)
	if err != nil{
		if errors.Is(err, repository.ErrUserNotFound){
			return nil, &models.Error{Code: models.USERNOTFOUND, Message: err}
		}
		return nil, &models.Error{Code: models.INTERNALERROR, Message: err}
	}

	return utypes, nil

}



func (s *UserService) GetAllUsers(ctx context.Context, limit, offset int32) ([]models.UserProfile, *models.Error) {
	users, err := s.userRepo.SelectAllUsers(ctx, limit, offset)
	if err != nil{
		return nil, &models.Error{Code: models.INTERNALERROR, Message: err}
	}
	return users, nil
}