package service

import (
	"auth_service/internal/models"
	"auth_service/internal/repository"
	"auth_service/pkg/kafka"
	"auth_service/pkg/token"
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	users       map[string]*models.User
	usersByEmail map[string]*models.User
	createErr   error
	getErr      error
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:       make(map[string]*models.User),
		usersByEmail: make(map[string]*models.User),
	}
}

func (m *mockUserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	if _, exists := m.usersByEmail[user.Email]; exists {
		return nil, repository.ErrUserExists
	}
	user.CreatedAt = time.Now()
	user.IsActive = true
	m.users[user.ID] = user
	m.usersByEmail[user.Email] = user
	return user, nil
}

func (m *mockUserRepository) DeleteUser(ctx context.Context, id string) error {
	if _, exists := m.users[id]; !exists {
		return repository.ErrUserNotFound
	}
	delete(m.users, id)
	return nil
}

func (m *mockUserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	user, exists := m.users[id]
	if !exists {
		return nil, repository.ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	user, exists := m.usersByEmail[email]
	if !exists {
		return nil, repository.ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserRepository) SetIsActive(ctx context.Context, id string, status bool) error {
	user, exists := m.users[id]
	if !exists {
		return repository.ErrUserNotFound
	}
	user.IsActive = status
	return nil
}

func (m *mockUserRepository) UpdatePassword(ctx context.Context, id string, password string) error {
	user, exists := m.users[id]
	if !exists {
		return repository.ErrUserNotFound
	}
	user.Password = password
	return nil
}

func (m *mockUserRepository) SelectAllUsers(ctx context.Context, limit, offset int32, isActive bool) ([]models.User, error) {
	var result []models.User
	for _, user := range m.users {
		if user.IsActive == isActive {
			result = append(result, *user)
		}
	}
	return result, nil
}

type mockTokenRepository struct {
	tokens   map[string]string
	saveErr  error
	getErr   error
	delErr   error
}

func newMockTokenRepository() *mockTokenRepository {
	return &mockTokenRepository{
		tokens: make(map[string]string),
	}
}

func (m *mockTokenRepository) SaveToken(ctx context.Context, userID, token string, expiresAt time.Time) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.tokens[token] = userID
	return nil
}

func (m *mockTokenRepository) GetToken(ctx context.Context, token string) (string, error) {
	if m.getErr != nil {
		return "", m.getErr
	}
	userID, exists := m.tokens[token]
	if !exists {
		return "", repository.ErrUserNotFound
	}
	return userID, nil
}

func (m *mockTokenRepository) DeleteToken(ctx context.Context, token string) error {
	if m.delErr != nil {
		return m.delErr
	}
	if _, exists := m.tokens[token]; !exists {
		return repository.ErrUserNotFound
	}
	delete(m.tokens, token)
	return nil
}

type mockProducer struct{}

func (m *mockProducer) Publish(ctx context.Context, topic, key string, value interface{}) error {
	return nil
}

func newTestAuthService(userRepo *mockUserRepository, tokenRepo *mockTokenRepository) *authService {
	tokenCfg := token.TokenConfig{
		Secret:     []byte("test-secret-key-for-testing-purposes"),
		AccessTTL:  5 * time.Minute,
		RefreshTTL: 24 * time.Hour,
	}

	producer := &kafka.Producer{
		Topic: "test-topic",
	}

	return &authService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		tokens:    token.NewManager(tokenCfg),
		producer:  producer,
	}
}

func TestRegister_Success(t *testing.T) {
	userRepo := newMockUserRepository()
	tokenRepo := newMockTokenRepository()
	svc := newTestAuthService(userRepo, tokenRepo)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"

	user, tokens, err := svc.Register(ctx, email, password, password)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if tokens == nil {
		t.Fatal("expected tokens, got nil")
	}
	if user.Email != email {
		t.Errorf("expected email %s, got %s", email, user.Email)
	}
	if tokens.AccessToken == "" {
		t.Error("expected access token, got empty")
	}
	if tokens.RefreshToken == "" {
		t.Error("expected refresh token, got empty")
	}
}

func TestRegister_PasswordMismatch(t *testing.T) {
	userRepo := newMockUserRepository()
	tokenRepo := newMockTokenRepository()
	svc := newTestAuthService(userRepo, tokenRepo)

	ctx := context.Background()
	email := "test@example.com"

	user, tokens, err := svc.Register(ctx, email, "password123", "different123")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != models.INVALIDINPUT {
		t.Errorf("expected error code %s, got %s", models.INVALIDINPUT, err.Code)
	}
	if user != nil {
		t.Error("expected nil user")
	}
	if tokens != nil {
		t.Error("expected nil tokens")
	}
}

func TestRegister_UserExists(t *testing.T) {
	userRepo := newMockUserRepository()
	tokenRepo := newMockTokenRepository()
	svc := newTestAuthService(userRepo, tokenRepo)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"

	_, _, _ = svc.Register(ctx, email, password, password)

	user, tokens, err := svc.Register(ctx, email, password, password)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != models.USEREXISTS {
		t.Errorf("expected error code %s, got %s", models.USEREXISTS, err.Code)
	}
	if user != nil {
		t.Error("expected nil user")
	}
	if tokens != nil {
		t.Error("expected nil tokens")
	}
}

func TestLogin_Success(t *testing.T) {
	userRepo := newMockUserRepository()
	tokenRepo := newMockTokenRepository()
	svc := newTestAuthService(userRepo, tokenRepo)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"

	_, _, _ = svc.Register(ctx, email, password, password)

	user, tokens, err := svc.Login(ctx, email, password)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if tokens == nil {
		t.Fatal("expected tokens, got nil")
	}
	if user.Email != email {
		t.Errorf("expected email %s, got %s", email, user.Email)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	userRepo := newMockUserRepository()
	tokenRepo := newMockTokenRepository()
	svc := newTestAuthService(userRepo, tokenRepo)

	ctx := context.Background()

	user, tokens, err := svc.Login(ctx, "nonexistent@example.com", "password123")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != models.INVALIDCREDENTIALS {
		t.Errorf("expected error code %s, got %s", models.INVALIDCREDENTIALS, err.Code)
	}
	if user != nil {
		t.Error("expected nil user")
	}
	if tokens != nil {
		t.Error("expected nil tokens")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	userRepo := newMockUserRepository()
	tokenRepo := newMockTokenRepository()
	svc := newTestAuthService(userRepo, tokenRepo)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"

	_, _, _ = svc.Register(ctx, email, password, password)

	user, tokens, err := svc.Login(ctx, email, "wrongpassword")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != models.INVALIDCREDENTIALS {
		t.Errorf("expected error code %s, got %s", models.INVALIDCREDENTIALS, err.Code)
	}
	if user != nil {
		t.Error("expected nil user")
	}
	if tokens != nil {
		t.Error("expected nil tokens")
	}
}

func TestLogout_Success(t *testing.T) {
	userRepo := newMockUserRepository()
	tokenRepo := newMockTokenRepository()
	svc := newTestAuthService(userRepo, tokenRepo)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"

	_, tokens, _ := svc.Register(ctx, email, password, password)

	err := svc.Logout(ctx, tokens.RefreshToken)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestLogout_InvalidToken(t *testing.T) {
	userRepo := newMockUserRepository()
	tokenRepo := newMockTokenRepository()
	svc := newTestAuthService(userRepo, tokenRepo)

	ctx := context.Background()

	err := svc.Logout(ctx, "invalid-token")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != models.INVALIDTOKEN {
		t.Errorf("expected error code %s, got %s", models.INVALIDTOKEN, err.Code)
	}
}

func TestRefresh_Success(t *testing.T) {
	userRepo := newMockUserRepository()
	tokenRepo := newMockTokenRepository()
	svc := newTestAuthService(userRepo, tokenRepo)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"

	_, tokens, _ := svc.Register(ctx, email, password, password)

	newTokens, err := svc.Refresh(ctx, tokens.RefreshToken)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if newTokens == nil {
		t.Fatal("expected tokens, got nil")
	}
	if newTokens.AccessToken == "" {
		t.Error("expected access token, got empty")
	}
	if newTokens.RefreshToken == "" {
		t.Error("expected refresh token, got empty")
	}
}

func TestRefresh_InvalidToken(t *testing.T) {
	userRepo := newMockUserRepository()
	tokenRepo := newMockTokenRepository()
	svc := newTestAuthService(userRepo, tokenRepo)

	ctx := context.Background()

	tokens, err := svc.Refresh(ctx, "invalid-token")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != models.INVALIDTOKEN {
		t.Errorf("expected error code %s, got %s", models.INVALIDTOKEN, err.Code)
	}
	if tokens != nil {
		t.Error("expected nil tokens")
	}
}

func TestValidateAccessToken_Success(t *testing.T) {
	userRepo := newMockUserRepository()
	tokenRepo := newMockTokenRepository()
	svc := newTestAuthService(userRepo, tokenRepo)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"

	user, tokens, _ := svc.Register(ctx, email, password, password)

	userID, exp, err := svc.ValidateAccessToken(ctx, tokens.AccessToken)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if userID != user.ID {
		t.Errorf("expected user ID %s, got %s", user.ID, userID)
	}
	if exp <= 0 {
		t.Error("expected positive expiration, got non-positive")
	}
}

func TestValidateAccessToken_InvalidToken(t *testing.T) {
	userRepo := newMockUserRepository()
	tokenRepo := newMockTokenRepository()
	svc := newTestAuthService(userRepo, tokenRepo)

	ctx := context.Background()

	userID, _, err := svc.ValidateAccessToken(ctx, "invalid-token")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != models.INVALIDTOKEN {
		t.Errorf("expected error code %s, got %s", models.INVALIDTOKEN, err.Code)
	}
	if userID != "" {
		t.Error("expected empty user ID")
	}
}

func TestPasswordHashing(t *testing.T) {
	password := "testpassword123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		t.Errorf("password comparison failed: %v", err)
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte("wrongpassword"))
	if err == nil {
		t.Error("expected error for wrong password")
	}
}

func TestMockUserRepository_CreateAndGet(t *testing.T) {
	repo := newMockUserRepository()
	ctx := context.Background()

	user := &models.User{
		ID:       "test-id",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	created, err := repo.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	if created.ID != user.ID {
		t.Errorf("expected ID %s, got %s", user.ID, created.ID)
	}

	retrieved, err := repo.GetUserByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}
	if retrieved.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, retrieved.Email)
	}
}

func TestMockUserRepository_DeleteUser(t *testing.T) {
	repo := newMockUserRepository()
	ctx := context.Background()

	user := &models.User{
		ID:       "test-id",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	_, _ = repo.CreateUser(ctx, user)

	err := repo.DeleteUser(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}

	_, err = repo.GetUserByID(ctx, user.ID)
	if !errors.Is(err, repository.ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestMockTokenRepository_SaveAndGet(t *testing.T) {
	repo := newMockTokenRepository()
	ctx := context.Background()

	token := "test-token"
	userID := "test-user-id"
	expiresAt := time.Now().Add(time.Hour)

	err := repo.SaveToken(ctx, userID, token, expiresAt)
	if err != nil {
		t.Fatalf("failed to save token: %v", err)
	}

	retrievedUserID, err := repo.GetToken(ctx, token)
	if err != nil {
		t.Fatalf("failed to get token: %v", err)
	}
	if retrievedUserID != userID {
		t.Errorf("expected user ID %s, got %s", userID, retrievedUserID)
	}
}

func TestMockTokenRepository_DeleteToken(t *testing.T) {
	repo := newMockTokenRepository()
	ctx := context.Background()

	token := "test-token"
	userID := "test-user-id"
	expiresAt := time.Now().Add(time.Hour)

	_ = repo.SaveToken(ctx, userID, token, expiresAt)

	err := repo.DeleteToken(ctx, token)
	if err != nil {
		t.Fatalf("failed to delete token: %v", err)
	}

	_, err = repo.GetToken(ctx, token)
	if !errors.Is(err, repository.ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}
