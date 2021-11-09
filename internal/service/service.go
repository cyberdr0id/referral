package service

import (
	"fmt"

	"github.com/cyberdr0id/referral/internal/repository"
	"github.com/cyberdr0id/referral/pkg/hash"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) SignUp(name, password string) (string, error) {
	pass, _ := hash.HashPassword(password)
	id, err := s.repo.CreateUser(name, pass)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Service) LogIn(name, password string) (string, error) {
	user, err := s.repo.GetUser(name)
	if err != nil {
		return "", err
	}

	if !hash.CheckPassowrdHash(password, user.Password) {
		return "", fmt.Errorf("wrong password for this user")
	}

	//TODO: generate tokens?

	return user.ID, nil
}

func (s *Service) SendCandidate(name, surname, fileID string) (string, error) {
	id, err := s.repo.AddCandidate(name, surname, fileID)
	if err != nil {
		return "", err
	}

	//TODO: create new request?

	return id, nil
}

func (s *Service) GetRequests(userID, filterType string) ([]repository.Request, error) {
	userRequests, err := s.repo.GetRequests(userID, filterType)
	if err != nil {
		return nil, err
	}

	return userRequests, nil
}

func (s *Service) DownloadCV(id string) (string, error) {
	id, err := s.repo.GetCVID(id)
	if err != nil {
		return "", err
	}

	// TODO: download cv from storage - storage

	return "example.com/path/to/file.extension", nil
}

func (s *Service) UpdateRequest(requestId, status string) error {
	err := s.repo.UpdateRequest(requestId, status)
	if err != nil {
		return err
	}

	return nil
}
