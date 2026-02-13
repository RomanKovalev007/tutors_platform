package token

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var(
	ErrTokenExpired = errors.New("token expired")
)

type TokenConfig struct {
	AccessTTL  time.Duration `env:"ACCESS_TTL_M" env-default:"300"` 
	RefreshTTL time.Duration `env:"REFRESH_TTL_H" env-default:"336"`
	Secret     []byte     `env:"SECRET,required"`   
}

type Manager struct {
	cfg TokenConfig
}


func NewManager(cfg TokenConfig) *Manager {
	return &Manager{cfg: cfg}
}


func (m *Manager) NewAccessToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": userID,
		"exp": time.Now().Add(m.cfg.AccessTTL).Unix(),
		"iat": time.Now().Unix(),
	})
	return token.SignedString(m.cfg.Secret)
}

func (m *Manager) NewRefreshToken() (string, time.Time, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", time.Time{}, err
	}

	token := base64.RawURLEncoding.EncodeToString(b)
	exp := time.Now().Add(m.cfg.RefreshTTL)

	return token, exp, nil
}


func (m *Manager) ValidateAccessToken(tokenStr string) (string, int64, error) {
	if tokenStr == "" {
		return "", 0, errors.New("token is empty")
	}

	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	token, err := parser.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.cfg.Secret, nil
	})

	if err != nil {
		return "", 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", 0, errors.New("token claims does not exist")
	}

	userID, ok := claims["id"].(string)
	if !ok || userID == "" {
		return "", 0, errors.New("token userID does not exist")
	}

	exp, ok := claims["exp"]
	if !ok {
		return "", 0, errors.New("exp claim is missing")
	}

	var expTime int64
	switch v := exp.(type) {
	case float64:
		expTime = int64(v)
	case int64:
		expTime = v
	default:
		return "", 0, errors.New("unsupported expiration type")
	}

	if time.Unix(expTime, 0).Before(time.Now()) {
		return userID, expTime, ErrTokenExpired
	}

	return userID, expTime, nil
}
