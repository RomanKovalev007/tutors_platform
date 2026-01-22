package service

import (
	"context"
	"errors"
	"user-service/internal/models"
	"user-service/internal/repository"
)

type StudentService struct {
	userRepo    UserProfileRepository
	studentRepo   StudentProfileRepository
}

func NewStudentService(
	userRepo UserProfileRepository,
	studentRepo StudentProfileRepository,
) *StudentService {
	return &StudentService{
		userRepo:    userRepo,
		studentRepo:   studentRepo,
	}
}


func (s *StudentService) CreateStudentProfile(ctx context.Context, student *models.StudentProfile) (*models.StudentProfile, *models.Error) {

	_, err := s.userRepo.GetUserByID(ctx, student.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound){
			return nil, &models.Error{Code: models.USERNOTFOUND, Message: err}
		}
		return nil, &models.Error{Code: models.INTERNALERROR, Message: err}
	}

	studentProfile, err := s.studentRepo.CreateStudentProfile(ctx, student)
	if err != nil{
		if errors.Is(err, repository.ErrUserExists){
			return nil, &models.Error{Code: models.TUTOREXISTS, Message: err}
		}
		return nil, &models.Error{Code: models.INTERNALERROR, Message: err}
	}

	return studentProfile, nil
}

func (s *StudentService) GetStudentProfile(ctx context.Context, id string) (*models.StudentProfile, *models.Error) {
	studentProfile, err := s.studentRepo.GetStudentProfileByID(ctx, id)
	if err != nil{
		if errors.Is(err, repository.ErrUserNotFound){
			return nil, &models.Error{Code: models.TUTORNOTFOUND, Message: err}
		}
		return nil, &models.Error{Code: models.INTERNALERROR, Message: err}
	}

	return studentProfile, nil
}

func (s *StudentService) UpdateStudentProfile(ctx context.Context, student *models.StudentProfile) (*models.StudentProfile, *models.Error) {
	studentProfile, err := s.studentRepo.UpdateStudentProfile(ctx, student)
	if err != nil{
		if errors.Is(err, repository.ErrUserNotFound){
			return nil, &models.Error{Code: models.TUTORNOTFOUND, Message: err}
		}
		return nil, &models.Error{Code: models.INTERNALERROR, Message: err}	
	}

	return studentProfile, nil 
}

func (s *StudentService) DeleteStudentProfile(ctx context.Context, id string) *models.Error {
	err := s.studentRepo.DeleteStudentProfie(ctx, id)
	if err != nil{
		if errors.Is(err, repository.ErrUserNotFound){
			return &models.Error{Code: models.TUTORNOTFOUND, Message: err}
		}
		return &models.Error{Code: models.INTERNALERROR, Message: err}	
	}

	return nil 
}