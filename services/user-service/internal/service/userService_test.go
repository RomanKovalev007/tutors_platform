package service

import (
	"context"
	"testing"
	"time"
	"user-service/internal/models"
	"user-service/internal/repository"
)

type mockUserProfileRepository struct {
	users        map[string]*models.UserProfile
	usersByEmail map[string]*models.UserProfile
	createErr    error
	getErr       error
	updateErr    error
	deleteErr    error
}

func newMockUserProfileRepository() *mockUserProfileRepository {
	return &mockUserProfileRepository{
		users:        make(map[string]*models.UserProfile),
		usersByEmail: make(map[string]*models.UserProfile),
	}
}

func (m *mockUserProfileRepository) CreateUser(ctx context.Context, user *models.UserProfile) (*models.UserProfile, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	if _, exists := m.usersByEmail[user.Email]; exists {
		return nil, repository.ErrUserExists
	}
	user.CreatedAt = time.Now()
	m.users[user.UserID] = user
	m.usersByEmail[user.Email] = user
	return user, nil
}

func (m *mockUserProfileRepository) DeleteUser(ctx context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	user, exists := m.users[id]
	if !exists {
		return repository.ErrUserNotFound
	}
	delete(m.usersByEmail, user.Email)
	delete(m.users, id)
	return nil
}

func (m *mockUserProfileRepository) GetUserByEmail(ctx context.Context, email string) (*models.UserProfile, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	user, exists := m.usersByEmail[email]
	if !exists {
		return nil, repository.ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserProfileRepository) GetUserByID(ctx context.Context, id string) (*models.UserProfile, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	user, exists := m.users[id]
	if !exists {
		return nil, repository.ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserProfileRepository) GetUserTypes(ctx context.Context, id string) (*models.UserType, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	user, exists := m.users[id]
	if !exists {
		return nil, repository.ErrUserNotFound
	}
	return &models.UserType{
		IsStudent: user.IsStudent,
		IsTutor:   user.IsTutor,
	}, nil
}

func (m *mockUserProfileRepository) UpdateUser(ctx context.Context, user *models.UserProfile) (*models.UserProfile, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	existing, exists := m.users[user.UserID]
	if !exists {
		return nil, repository.ErrUserNotFound
	}
	if user.Name != "" {
		existing.Name = user.Name
	}
	if user.Surname != "" {
		existing.Surname = user.Surname
	}
	if user.Telegram != "" {
		existing.Telegram = user.Telegram
	}
	return existing, nil
}

func (m *mockUserProfileRepository) SelectAllUsers(ctx context.Context, limit, offset int32) ([]models.UserProfile, error) {
	var result []models.UserProfile
	for _, user := range m.users {
		result = append(result, *user)
	}
	return result, nil
}

type mockTutorProfileRepository struct {
	tutors    map[string]*models.TutorProfile
	createErr error
	getErr    error
	updateErr error
	deleteErr error
}

func newMockTutorProfileRepository() *mockTutorProfileRepository {
	return &mockTutorProfileRepository{
		tutors: make(map[string]*models.TutorProfile),
	}
}

func (m *mockTutorProfileRepository) CreateTutorProfile(ctx context.Context, tutor *models.TutorProfile) (*models.TutorProfile, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	m.tutors[tutor.UserID] = tutor
	return tutor, nil
}

func (m *mockTutorProfileRepository) DeleteTutorProfie(ctx context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.tutors, id)
	return nil
}

func (m *mockTutorProfileRepository) GetTutorProfileByID(ctx context.Context, id string) (*models.TutorProfile, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	tutor, exists := m.tutors[id]
	if !exists {
		return nil, repository.ErrUserNotFound
	}
	return tutor, nil
}

func (m *mockTutorProfileRepository) UpdateTutorProfile(ctx context.Context, tutor *models.TutorProfile) (*models.TutorProfile, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	existing, exists := m.tutors[tutor.UserID]
	if !exists {
		return nil, repository.ErrUserNotFound
	}
	if tutor.Specialization != "" {
		existing.Specialization = tutor.Specialization
	}
	if tutor.Bio != "" {
		existing.Bio = tutor.Bio
	}
	return existing, nil
}

type mockStudentProfileRepository struct {
	students  map[string]*models.StudentProfile
	createErr error
	getErr    error
	updateErr error
	deleteErr error
}

func newMockStudentProfileRepository() *mockStudentProfileRepository {
	return &mockStudentProfileRepository{
		students: make(map[string]*models.StudentProfile),
	}
}

func (m *mockStudentProfileRepository) CreateStudentProfile(ctx context.Context, student *models.StudentProfile) (*models.StudentProfile, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	m.students[student.UserID] = student
	return student, nil
}

func (m *mockStudentProfileRepository) DeleteStudentProfie(ctx context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.students, id)
	return nil
}

func (m *mockStudentProfileRepository) GetStudentProfileByID(ctx context.Context, id string) (*models.StudentProfile, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	student, exists := m.students[id]
	if !exists {
		return nil, repository.ErrUserNotFound
	}
	return student, nil
}

func (m *mockStudentProfileRepository) UpdateStudentProfile(ctx context.Context, student *models.StudentProfile) (*models.StudentProfile, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	existing, exists := m.students[student.UserID]
	if !exists {
		return nil, repository.ErrUserNotFound
	}
	if student.Bio != "" {
		existing.Bio = student.Bio
	}
	return existing, nil
}

func newTestUserService() (*UserService, *mockUserProfileRepository, *mockTutorProfileRepository, *mockStudentProfileRepository) {
	userRepo := newMockUserProfileRepository()
	tutorRepo := newMockTutorProfileRepository()
	studentRepo := newMockStudentProfileRepository()
	svc := NewUserService(userRepo, tutorRepo, studentRepo)
	return svc, userRepo, tutorRepo, studentRepo
}

func TestUserService_CreateUser_Success(t *testing.T) {
	svc, _, _, _ := newTestUserService()
	ctx := context.Background()

	user := &models.UserProfile{
		UserID:  "test-id",
		Email:   "test@example.com",
		Name:    "Test",
		Surname: "User",
	}

	result, err := svc.CreateUser(ctx, user)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected user, got nil")
	}
	if result.UserID != user.UserID {
		t.Errorf("expected ID %s, got %s", user.UserID, result.UserID)
	}
}

func TestUserService_CreateUser_AlreadyExists(t *testing.T) {
	svc, _, _, _ := newTestUserService()
	ctx := context.Background()

	user := &models.UserProfile{
		UserID:  "test-id",
		Email:   "test@example.com",
		Name:    "Test",
		Surname: "User",
	}

	_, _ = svc.CreateUser(ctx, user)

	user2 := &models.UserProfile{
		UserID:  "test-id-2",
		Email:   "test@example.com",
		Name:    "Test2",
		Surname: "User2",
	}

	result, err := svc.CreateUser(ctx, user2)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != models.USEREXISTS {
		t.Errorf("expected error code %s, got %s", models.USEREXISTS, err.Code)
	}
	if result != nil {
		t.Error("expected nil user")
	}
}

func TestUserService_GetUserByID_Success(t *testing.T) {
	svc, _, _, _ := newTestUserService()
	ctx := context.Background()

	user := &models.UserProfile{
		UserID:  "test-id",
		Email:   "test@example.com",
		Name:    "Test",
		Surname: "User",
	}
	_, _ = svc.CreateUser(ctx, user)

	result, err := svc.GetUserByID(ctx, user.UserID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected user, got nil")
	}
	if result.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, result.Email)
	}
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	svc, _, _, _ := newTestUserService()
	ctx := context.Background()

	result, err := svc.GetUserByID(ctx, "nonexistent-id")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != models.USERNOTFOUND {
		t.Errorf("expected error code %s, got %s", models.USERNOTFOUND, err.Code)
	}
	if result != nil {
		t.Error("expected nil user")
	}
}

func TestUserService_GetUserByEmail_Success(t *testing.T) {
	svc, _, _, _ := newTestUserService()
	ctx := context.Background()

	user := &models.UserProfile{
		UserID:  "test-id",
		Email:   "test@example.com",
		Name:    "Test",
		Surname: "User",
	}
	_, _ = svc.CreateUser(ctx, user)

	result, err := svc.GetUserByEmail(ctx, user.Email)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected user, got nil")
	}
	if result.UserID != user.UserID {
		t.Errorf("expected ID %s, got %s", user.UserID, result.UserID)
	}
}

func TestUserService_GetUserByEmail_NotFound(t *testing.T) {
	svc, _, _, _ := newTestUserService()
	ctx := context.Background()

	result, err := svc.GetUserByEmail(ctx, "nonexistent@example.com")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != models.USERNOTFOUND {
		t.Errorf("expected error code %s, got %s", models.USERNOTFOUND, err.Code)
	}
	if result != nil {
		t.Error("expected nil user")
	}
}

func TestUserService_UpdateUser_Success(t *testing.T) {
	svc, _, _, _ := newTestUserService()
	ctx := context.Background()

	user := &models.UserProfile{
		UserID:  "test-id",
		Email:   "test@example.com",
		Name:    "Test",
		Surname: "User",
	}
	_, _ = svc.CreateUser(ctx, user)

	update := &models.UserProfile{
		UserID: "test-id",
		Name:   "Updated",
	}

	result, err := svc.UpdateUser(ctx, update)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected user, got nil")
	}
	if result.Name != "Updated" {
		t.Errorf("expected name Updated, got %s", result.Name)
	}
}

func TestUserService_UpdateUser_NotFound(t *testing.T) {
	svc, _, _, _ := newTestUserService()
	ctx := context.Background()

	update := &models.UserProfile{
		UserID: "nonexistent-id",
		Name:   "Updated",
	}

	result, err := svc.UpdateUser(ctx, update)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != models.USERNOTFOUND {
		t.Errorf("expected error code %s, got %s", models.USERNOTFOUND, err.Code)
	}
	if result != nil {
		t.Error("expected nil user")
	}
}

func TestUserService_DeleteUser_Success(t *testing.T) {
	svc, _, _, _ := newTestUserService()
	ctx := context.Background()

	user := &models.UserProfile{
		UserID:  "test-id",
		Email:   "test@example.com",
		Name:    "Test",
		Surname: "User",
	}
	_, _ = svc.CreateUser(ctx, user)

	err := svc.DeleteUser(ctx, user.UserID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, getErr := svc.GetUserByID(ctx, user.UserID)
	if getErr == nil {
		t.Fatal("expected error after deletion")
	}
}

func TestUserService_DeleteUser_NotFound(t *testing.T) {
	svc, _, _, _ := newTestUserService()
	ctx := context.Background()

	err := svc.DeleteUser(ctx, "nonexistent-id")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != models.USERNOTFOUND {
		t.Errorf("expected error code %s, got %s", models.USERNOTFOUND, err.Code)
	}
}

func TestUserService_DeleteUser_WithTutorProfile(t *testing.T) {
	svc, userRepo, tutorRepo, _ := newTestUserService()
	ctx := context.Background()

	user := &models.UserProfile{
		UserID:  "test-id",
		Email:   "test@example.com",
		Name:    "Test",
		Surname: "User",
		UserType: models.UserType{
			IsTutor: true,
		},
	}
	_, _ = userRepo.CreateUser(ctx, user)

	tutor := &models.TutorProfile{
		UserID:         "test-id",
		Specialization: "Math",
	}
	_, _ = tutorRepo.CreateTutorProfile(ctx, tutor)

	err := svc.DeleteUser(ctx, user.UserID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUserService_GetUserTypes_Success(t *testing.T) {
	svc, _, _, _ := newTestUserService()
	ctx := context.Background()

	user := &models.UserProfile{
		UserID:  "test-id",
		Email:   "test@example.com",
		Name:    "Test",
		Surname: "User",
		UserType: models.UserType{
			IsTutor:   true,
			IsStudent: false,
		},
	}
	_, _ = svc.CreateUser(ctx, user)

	result, err := svc.GetUserTypes(ctx, user.UserID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected user types, got nil")
	}
	if !result.IsTutor {
		t.Error("expected IsTutor to be true")
	}
	if result.IsStudent {
		t.Error("expected IsStudent to be false")
	}
}

func TestUserService_GetUserTypes_NotFound(t *testing.T) {
	svc, _, _, _ := newTestUserService()
	ctx := context.Background()

	result, err := svc.GetUserTypes(ctx, "nonexistent-id")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != models.USERNOTFOUND {
		t.Errorf("expected error code %s, got %s", models.USERNOTFOUND, err.Code)
	}
	if result != nil {
		t.Error("expected nil user types")
	}
}

func TestUserService_GetAllUsers_Success(t *testing.T) {
	svc, _, _, _ := newTestUserService()
	ctx := context.Background()

	user1 := &models.UserProfile{UserID: "id1", Email: "test1@example.com", Name: "User1"}
	user2 := &models.UserProfile{UserID: "id2", Email: "test2@example.com", Name: "User2"}
	_, _ = svc.CreateUser(ctx, user1)
	_, _ = svc.CreateUser(ctx, user2)

	users, err := svc.GetAllUsers(ctx, 10, 0)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestUserService_GetAllUsers_Empty(t *testing.T) {
	svc, _, _, _ := newTestUserService()
	ctx := context.Background()

	users, err := svc.GetAllUsers(ctx, 10, 0)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if users == nil {
		users = []models.UserProfile{}
	}
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}
}
