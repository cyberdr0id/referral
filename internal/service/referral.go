package service

import (
	"errors"

	"github.com/cyberdr0id/referral/internal/repository"
)

type ReferralService struct {
	repo *repository.Repository
}

func NewReferralService(repo *repository.Repository) *ReferralService {
	return &ReferralService{
		repo: repo,
	}
}

func (s *ReferralService) SendCandidate(userID, name, surname, fileID string) (string, error) {
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

func (s *ReferralService) GetRequests(userID, filterType string) ([]repository.Request, error) {
	userRequests, err := s.repo.GetRequests(userID, filterType)
	if err != nil {
		return nil, err
	}

	return userRequests, nil
}

func (s *ReferralService) DownloadCV(id string) (string, error) {
	_, err := s.repo.GetCVID(id)
	if err != nil {
		return "", err
	}

	// TODO: download cv from storage - storage

	return "example.com/path/to/file.extension", nil
}

func (s *ReferralService) UpdateRequest(userID, requestId, status string) error {
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
