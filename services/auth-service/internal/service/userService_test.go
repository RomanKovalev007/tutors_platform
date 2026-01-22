package service

import (
	"auth_service/internal/models"
	"context"
	"testing"
)

func newTestUserService(userRepo *mockUserRepository) *userService {
	return &userService{
		userRepo: userRepo,
	}
}

func TestUserService_GetUserByID_Success(t *testing.T) {
	userRepo := newMockUserRepository()
	svc := newTestUserService(userRepo)
	ctx := context.Background()

	user := &models.User{
		ID:       "test-id",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	_, _ = userRepo.CreateUser(ctx, user)

	result, err := svc.GetUserByID(ctx, user.ID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected user, got nil")
	}
	if result.ID != user.ID {
		t.Errorf("expected ID %s, got %s", user.ID, result.ID)
	}
	if result.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, result.Email)
	}
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	userRepo := newMockUserRepository()
	svc := newTestUserService(userRepo)
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
	userRepo := newMockUserRepository()
	svc := newTestUserService(userRepo)
	ctx := context.Background()

	user := &models.User{
		ID:       "test-id",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	_, _ = userRepo.CreateUser(ctx, user)

	result, err := svc.GetUserByEmail(ctx, user.Email)

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

func TestUserService_GetUserByEmail_NotFound(t *testing.T) {
	userRepo := newMockUserRepository()
	svc := newTestUserService(userRepo)
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

func TestUserService_UpdateIsActiveStatus_Success(t *testing.T) {
	userRepo := newMockUserRepository()
	svc := newTestUserService(userRepo)
	ctx := context.Background()

	user := &models.User{
		ID:       "test-id",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	_, _ = userRepo.CreateUser(ctx, user)

	err := svc.UpdateIsActiveStatus(ctx, user.ID, false)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, _ := userRepo.GetUserByID(ctx, user.ID)
	if updated.IsActive {
		t.Error("expected user to be inactive")
	}
}

func TestUserService_UpdateIsActiveStatus_NotFound(t *testing.T) {
	userRepo := newMockUserRepository()
	svc := newTestUserService(userRepo)
	ctx := context.Background()

	err := svc.UpdateIsActiveStatus(ctx, "nonexistent-id", false)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != models.USERNOTFOUND {
		t.Errorf("expected error code %s, got %s", models.USERNOTFOUND, err.Code)
	}
}

func TestUserService_DeleteUser_Success(t *testing.T) {
	userRepo := newMockUserRepository()
	svc := newTestUserService(userRepo)
	ctx := context.Background()

	user := &models.User{
		ID:       "test-id",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	_, _ = userRepo.CreateUser(ctx, user)

	err := svc.DeleteUser(ctx, user.ID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, getErr := svc.GetUserByID(ctx, user.ID)
	if getErr == nil {
		t.Fatal("expected error after deletion")
	}
	if getErr.Code != models.USERNOTFOUND {
		t.Errorf("expected error code %s, got %s", models.USERNOTFOUND, getErr.Code)
	}
}

func TestUserService_DeleteUser_NotFound(t *testing.T) {
	userRepo := newMockUserRepository()
	svc := newTestUserService(userRepo)
	ctx := context.Background()

	err := svc.DeleteUser(ctx, "nonexistent-id")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != models.USERNOTFOUND {
		t.Errorf("expected error code %s, got %s", models.USERNOTFOUND, err.Code)
	}
}

func TestUserService_GetAllUsers_Success(t *testing.T) {
	userRepo := newMockUserRepository()
	svc := newTestUserService(userRepo)
	ctx := context.Background()

	user1 := &models.User{ID: "id1", Email: "test1@example.com", Password: "hash1"}
	user2 := &models.User{ID: "id2", Email: "test2@example.com", Password: "hash2"}
	_, _ = userRepo.CreateUser(ctx, user1)
	_, _ = userRepo.CreateUser(ctx, user2)

	users, err := svc.GetAllUsers(ctx, 10, 0, true)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestUserService_GetAllUsers_Empty(t *testing.T) {
	userRepo := newMockUserRepository()
	svc := newTestUserService(userRepo)
	ctx := context.Background()

	users, err := svc.GetAllUsers(ctx, 10, 0, true)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if users == nil {
		users = []models.User{}
	}
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}
}
