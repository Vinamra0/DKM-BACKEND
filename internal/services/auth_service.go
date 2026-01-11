package services

import (
	"errors"
	"time"

	"dkmbackend/internal/config"
	"dkmbackend/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct{ cfg config.Config }

func NewAuthService(cfg config.Config) *AuthService { return &AuthService{cfg: cfg} }

func (s *AuthService) Login(email, password string) (string, *models.User, error) {
	if email != s.cfg.AdminEmail || password != s.cfg.AdminPassword {
		return "", nil, errors.New("invalid credentials")
	}
	user := &models.User{Email: s.cfg.AdminEmail, Name: s.cfg.AdminName}
	token, err := s.generateToken(user)
	return token, user, err
}

func (s *AuthService) generateToken(u *models.User) (string, error) {
	claims := jwt.MapClaims{
		"email": u.Email,
		"name":  u.Name,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *AuthService) ParseToken(tokenStr string) (*models.User, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		email, _ := claims["email"].(string)
		name, _ := claims["name"].(string)
		return &models.User{Email: email, Name: name}, nil
	}
	return nil, errors.New("invalid token claims")
}
