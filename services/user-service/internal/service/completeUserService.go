package service

import (
	"context"
	"errors"
	"user-service/internal/models"
	"user-service/internal/repository"
)

func (s *UserService) GetCompleteUserProfile(ctx context.Context, id string) (*models.CompleteUserProfile, *models.Error) {

	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, &models.Error{Code: models.USERNOTFOUND, Message: err}
		}
		return nil, &models.Error{Code: models.INTERNALERROR, Message: err}
	}

	completeProfile := &models.CompleteUserProfile{
		UserProfile: user,
	}

	if user.IsTutor {
		tutorProfile, err := s.tutorRepo.GetTutorProfileByID(ctx, id)
		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				return nil, &models.Error{Code: models.USERNOTFOUND, Message: err}
			}
			return nil, &models.Error{Code: models.INTERNALERROR, Message: err}
		}
		completeProfile.TutorProfile = tutorProfile
	}

	if user.IsStudent {
		studentProfile, err := s.studentRepo.GetStudentProfileByID(ctx, id)
		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				return nil, &models.Error{Code: models.USERNOTFOUND, Message: err}
			}
			return nil, &models.Error{Code: models.INTERNALERROR, Message: err}
		}
		completeProfile.StudentProfile = studentProfile
	}

	return completeProfile, nil
}