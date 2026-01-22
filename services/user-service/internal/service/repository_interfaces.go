package service

import (
	"context"
	"user-service/internal/models"
)

type UserProfileRepository interface {
	CreateUser(ctx context.Context, user *models.UserProfile) (*models.UserProfile, error)
	DeleteUser(ctx context.Context, id string) error
	GetUserByEmail(ctx context.Context, email string) (*models.UserProfile, error)
	GetUserByID(ctx context.Context, id string) (*models.UserProfile, error)
	GetUserTypes(ctx context.Context, id string) (*models.UserType, error)
	UpdateUser(ctx context.Context, user *models.UserProfile) (*models.UserProfile, error)
	SelectAllUsers(ctx context.Context, limit, offset int32) ([]models.UserProfile, error)
}

type TutorProfileRepository interface {
	CreateTutorProfile(ctx context.Context, tutor *models.TutorProfile) (*models.TutorProfile, error)
	DeleteTutorProfie(ctx context.Context, id string) error
	GetTutorProfileByID(ctx context.Context, id string) (*models.TutorProfile, error)
	UpdateTutorProfile(ctx context.Context, tutor *models.TutorProfile) (*models.TutorProfile, error)
}

type StudentProfileRepository interface {
	CreateStudentProfile(ctx context.Context, student *models.StudentProfile) (*models.StudentProfile, error)
	DeleteStudentProfie(ctx context.Context, id string) error
	GetStudentProfileByID(ctx context.Context, id string) (*models.StudentProfile, error)
	UpdateStudentProfile(ctx context.Context, student *models.StudentProfile) (*models.StudentProfile, error)
}
