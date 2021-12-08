package service

import (
	"errors"
	"fmt"

	"github.com/cyberdr0id/referral/internal/repository"
	"github.com/cyberdr0id/referral/pkg/hash"
	"github.com/cyberdr0id/referral/pkg/jwt"
)

// AuthService present a service for authorization service.
type AuthService struct {
	repo         *repository.Repository
	tokenManager *jwt.TokenManager
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(repo *repository.Repository, tm *jwt.TokenManager) *AuthService {
	return &AuthService{
		repo:         repo,
		tokenManager: tm,
	}
}

// SignUp hash password and add user to database.
func (s *AuthService) SignUp(name, password string) (string, error) {
	pass, err := hash.HashPassword(password)
	if err != nil {
		return "", fmt.Errorf("unable to hash password: %w", err)
	}

	id, err := s.repo.CreateUser(name, pass)
	if err != nil {
		return "", fmt.Errorf("cannot create user: %w", err)
	}

	return id, nil
}

// LogIn gets user from database, comparing passwords and generate JWT token - auathorize user.
func (s *AuthService) LogIn(name, password string) (string, error) {
	user, err := s.repo.GetUser(name)
	if errors.Is(err, repository.ErrNoUser) {
		return "", ErrNoUser
	}
	if err != nil {
		return "", fmt.Errorf("cannot get user from database: %w", err)
	}

	ok := hash.CheckPasswordHash(password, user.Password)
	if !ok {
		return "", ErrNoUser
	}

	token, err := s.tokenManager.GenerateToken(user.ID, user.IsAdmin)
	if err != nil {
		return "", fmt.Errorf("cannot generate JWT token: %w", err)
	}

	return token, nil
}
