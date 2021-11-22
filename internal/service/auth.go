package service

import (
	"fmt"

	"github.com/cyberdr0id/referral/internal/repository"
	"github.com/cyberdr0id/referral/pkg/hash"
)

type AuthService struct {
	repo *repository.Repository
}

func NewAuthService(repo *repository.Repository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) SignUp(name, password string) (string, error) {
	pass, _ := hash.HashPassword(password)
	id, err := s.repo.CreateUser(name, pass)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *AuthService) LogIn(name, password string) (string, error) {
	user, err := s.repo.GetUser(name)
	if err != nil {
		return "", err
	}

	if !hash.CheckPassowrdHash(password, user.Password) {
		return "", fmt.Errorf("wrong password for this user") // TODO: is this good practice?
	}

	//TODO: generate tokens?

	return user.ID, nil
}
