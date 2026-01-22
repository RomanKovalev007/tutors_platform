package token

import (
	"testing"
	"time"
)

func newTestManager() *Manager {
	cfg := TokenConfig{
		Secret:     []byte("test-secret-key-for-testing"),
		AccessTTL:  5 * time.Minute,
		RefreshTTL: 24 * time.Hour,
	}
	return NewManager(cfg)
}

func TestNewAccessToken_Success(t *testing.T) {
	m := newTestManager()
	userID := "test-user-id"

	token, err := m.NewAccessToken(userID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
}

func TestNewRefreshToken_Success(t *testing.T) {
	m := newTestManager()

	token, exp, err := m.NewRefreshToken()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
	if exp.Before(time.Now()) {
		t.Error("expected future expiration time")
	}
}

func TestValidateAccessToken_Success(t *testing.T) {
	m := newTestManager()
	userID := "test-user-id"

	token, err := m.NewAccessToken(userID)
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	validatedUserID, exp, err := m.ValidateAccessToken(token)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if validatedUserID != userID {
		t.Errorf("expected user ID %s, got %s", userID, validatedUserID)
	}
	if exp <= 0 {
		t.Error("expected positive expiration")
	}
}

func TestValidateAccessToken_InvalidToken(t *testing.T) {
	m := newTestManager()

	_, _, err := m.ValidateAccessToken("invalid-token")

	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestValidateAccessToken_EmptyToken(t *testing.T) {
	m := newTestManager()

	_, _, err := m.ValidateAccessToken("")

	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestValidateAccessToken_ExpiredToken(t *testing.T) {
	cfg := TokenConfig{
		Secret:     []byte("test-secret-key-for-testing"),
		AccessTTL:  -1 * time.Second,
		RefreshTTL: 24 * time.Hour,
	}
	m := NewManager(cfg)
	userID := "test-user-id"

	token, err := m.NewAccessToken(userID)
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	validatedUserID, _, err := m.ValidateAccessToken(token)

	if err == nil {
		t.Fatal("expected error for expired token")
	}
	if err != ErrTokenExpired {
		t.Errorf("expected ErrTokenExpired, got %v", err)
	}
	if validatedUserID != userID {
		t.Errorf("expected user ID %s even for expired token, got %s", userID, validatedUserID)
	}
}

func TestValidateAccessToken_WrongSecret(t *testing.T) {
	m1 := newTestManager()

	cfg2 := TokenConfig{
		Secret:     []byte("different-secret-key"),
		AccessTTL:  5 * time.Minute,
		RefreshTTL: 24 * time.Hour,
	}
	m2 := NewManager(cfg2)

	token, _ := m1.NewAccessToken("test-user-id")

	_, _, err := m2.ValidateAccessToken(token)

	if err == nil {
		t.Fatal("expected error for token with wrong secret")
	}
}

func TestNewRefreshToken_Uniqueness(t *testing.T) {
	m := newTestManager()

	tokens := make(map[string]bool)
	for i := 0; i < 100; i++ {
		token, _, err := m.NewRefreshToken()
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}
		if tokens[token] {
			t.Fatal("generated duplicate refresh token")
		}
		tokens[token] = true
	}
}

func TestNewAccessToken_DifferentUsers(t *testing.T) {
	m := newTestManager()

	token1, _ := m.NewAccessToken("user1")
	token2, _ := m.NewAccessToken("user2")

	if token1 == token2 {
		t.Error("expected different tokens for different users")
	}

	userID1, _, _ := m.ValidateAccessToken(token1)
	userID2, _, _ := m.ValidateAccessToken(token2)

	if userID1 != "user1" {
		t.Errorf("expected user1, got %s", userID1)
	}
	if userID2 != "user2" {
		t.Errorf("expected user2, got %s", userID2)
	}
}

func TestRefreshToken_ExpirationTime(t *testing.T) {
	cfg := TokenConfig{
		Secret:     []byte("test-secret"),
		AccessTTL:  5 * time.Minute,
		RefreshTTL: 1 * time.Hour,
	}
	m := NewManager(cfg)

	_, exp, err := m.NewRefreshToken()
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	expectedExp := time.Now().Add(1 * time.Hour)
	tolerance := 5 * time.Second

	if exp.Before(expectedExp.Add(-tolerance)) || exp.After(expectedExp.Add(tolerance)) {
		t.Errorf("expiration time %v not within expected range around %v", exp, expectedExp)
	}
}
