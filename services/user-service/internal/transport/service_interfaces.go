package transport

import (
	"context"
	"user-service/internal/models"
)

type UserService interface{
	CreateUser(ctx context.Context, user *models.UserProfile) (*models.UserProfile, *models.Error)
	DeleteUser(ctx context.Context, id string) *models.Error
	GetAllUsers(ctx context.Context, limit int32, offset int32) ([]models.UserProfile, *models.Error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserProfile, *models.Error)
	GetUserByID(ctx context.Context, id string) (*models.UserProfile, *models.Error)
	GetUserTypes(ctx context.Context, id string) (*models.UserType, *models.Error)
	UpdateUser(ctx context.Context, user *models.UserProfile) (*models.UserProfile, *models.Error)
	GetCompleteUserProfile(ctx context.Context, id string) (*models.CompleteUserProfile, *models.Error)
}

type TutorService interface {
	CreateTutorProfile(ctx context.Context, tutor *models.TutorProfile) (*models.TutorProfile, *models.Error)
	DeleteTutorProfile(ctx context.Context, id string) *models.Error
	GetTutorProfile(ctx context.Context, id string) (*models.TutorProfile, *models.Error)
	UpdateTutorProfile(ctx context.Context, tutor *models.TutorProfile) (*models.TutorProfile, *models.Error)
	ValidateTutor(ctx context.Context, id string) (bool, *models.Error)
}

type StudentService interface {
	CreateStudentProfile(ctx context.Context, student *models.StudentProfile) (*models.StudentProfile, *models.Error)
	DeleteStudentProfile(ctx context.Context, id string) *models.Error
	GetStudentProfile(ctx context.Context, id string) (*models.StudentProfile, *models.Error)
	UpdateStudentProfile(ctx context.Context, student *models.StudentProfile) (*models.StudentProfile, *models.Error)
}