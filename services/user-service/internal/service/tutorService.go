package service

import (
	"context"
	"errors"
	"user-service/internal/models"
	"user-service/internal/repository"
)

type TutorService struct {
	userRepo    UserProfileRepository
	tutorRepo   TutorProfileRepository
}

func NewTutorService(
	userRepo UserProfileRepository,
	tutorRepo TutorProfileRepository,
) *TutorService {
	return &TutorService{
		userRepo:    userRepo,
		tutorRepo:   tutorRepo,
	}
}

func (s *TutorService) CreateTutorProfile(ctx context.Context, tutor *models.TutorProfile) (*models.TutorProfile, *models.Error) {

	_, err := s.userRepo.GetUserByID(ctx, tutor.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound){
			return nil, &models.Error{Code: models.USERNOTFOUND, Message: err}
		}
		return nil, &models.Error{Code: models.INTERNALERROR, Message: err}
	}

	tutorProfile, err := s.tutorRepo.CreateTutorProfile(ctx, tutor)
	if err != nil{
		if errors.Is(err, repository.ErrUserExists){
			return nil, &models.Error{Code: models.TUTOREXISTS, Message: err}
		}
		return nil, &models.Error{Code: models.INTERNALERROR, Message: err}
	}

	return tutorProfile, nil
}

func (s *TutorService) GetTutorProfile(ctx context.Context, id string) (*models.TutorProfile, *models.Error) {
	tutorProfile, err := s.tutorRepo.GetTutorProfileByID(ctx, id)
	if err != nil{
		if errors.Is(err, repository.ErrUserNotFound){
			return nil, &models.Error{Code: models.TUTORNOTFOUND, Message: err}
		}
		return nil, &models.Error{Code: models.INTERNALERROR, Message: err}
	}

	return tutorProfile, nil
}

func (s *TutorService) UpdateTutorProfile(ctx context.Context, tutor *models.TutorProfile) (*models.TutorProfile, *models.Error) {
	tutorProfile, err := s.tutorRepo.UpdateTutorProfile(ctx, tutor)
	if err != nil{
		if errors.Is(err, repository.ErrUserNotFound){
			return nil, &models.Error{Code: models.TUTORNOTFOUND, Message: err}
		}
		return nil, &models.Error{Code: models.INTERNALERROR, Message: err}	
	}

	return tutorProfile, nil 
}

func (s *TutorService) DeleteTutorProfile(ctx context.Context, id string) *models.Error {
	err := s.tutorRepo.DeleteTutorProfie(ctx, id)
	if err != nil{
		if errors.Is(err, repository.ErrUserNotFound){
			return &models.Error{Code: models.TUTORNOTFOUND, Message: err}
		}
		return &models.Error{Code: models.INTERNALERROR, Message: err}	
	}

	return nil 
}

func (s *TutorService) ValidateTutor(ctx context.Context, id string) (bool, *models.Error){
	types, err := s.userRepo.GetUserTypes(ctx, id)
	if err != nil{
		return false, &models.Error{Code: models.INTERNALERROR, Message: err}
	}

	return types.IsTutor, nil
}