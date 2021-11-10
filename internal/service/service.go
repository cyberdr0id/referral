package service

import (
	"errors"
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
		return "", fmt.Errorf("wrong password for this user") // TODO: is this good practice?
	}

	//TODO: generate tokens?

	return user.ID, nil
}

func (s *Service) SendCandidate(userID, name, surname, fileID string) (string, error) {
	candidateID, err := s.repo.AddCandidate(name, surname, fileID)
	if err != nil {
		return "", err
	}

	id, err := s.repo.CreateRequest(userID, candidateID)
	if err != nil {
		return "", nil
	}

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
	_, err := s.repo.GetCVID(id)
	if err != nil {
		return "", err
	}

	// TODO: download cv from storage - storage

	return "example.com/path/to/file.extension", nil
}

func (s *Service) UpdateRequest(userID, requestId, status string) error {
	isAdmin, err := s.repo.IsUserAdmin(userID)
	if err != nil {
		return err
	}

	err = s.repo.IsUserRequest(userID, requestId)
	if errors.Is(err, repository.ErrNoAccess) && !isAdmin {
		return err
	}
	if err != nil {
		return err
	}

	if err := s.repo.UpdateRequest(requestId, status); err != nil {
		return err
	}

	return nil
}
