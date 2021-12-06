package service

import (
	"fmt"

	"github.com/cyberdr0id/referral/internal/repository"
	"github.com/cyberdr0id/referral/pkg/hash"
	"github.com/cyberdr0id/referral/pkg/jwt"
)

type AuthService struct {
	repo         *repository.Repository
	tokenManager *jwt.TokenManager
}

func NewAuthService(repo *repository.Repository, tm *jwt.TokenManager) *AuthService {
	return &AuthService{
		repo:         repo,
		tokenManager: tm,
	}
}

func (s *AuthService) CreateUser(name, password string) (string, error) {
	pass, err := hash.HashPassword(password)
	if err != nil {
		return "", fmt.Errorf("unable to hash password: %w", err)
	}

	id, err := s.repo.CreateUser(name, pass)
	if err != nil {
		return "", err
	}

	return id, nil
}
