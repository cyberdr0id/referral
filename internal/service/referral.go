package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	mycontext "github.com/cyberdr0id/referral/internal/context"
	"github.com/cyberdr0id/referral/internal/repository"
	"github.com/cyberdr0id/referral/internal/storage"
	"github.com/pborman/uuid"
)

// ReferralService presents access to referral service via repository.
type ReferralService struct {
	repo *repository.Repository
	gcs  *storage.Storage
}

// NewReferralService creates a new instance of ReferralService.
func NewReferralService(repo *repository.Repository, gcs *storage.Storage) *ReferralService {
	return &ReferralService{
		repo: repo,
		gcs:  gcs,
	}
}

// SubmitCandidateRequest presents a type for reading data after submitting a candidate.
type SubmitCandidateRequest struct {
	File             multipart.File
	CandidateName    string
	CandidateSurname string
}

// AddCandidate creates request with candidate.
func (s *ReferralService) AddCandidate(ctx context.Context, request SubmitCandidateRequest) (string, error) {
	userID, ok := mycontext.GetUserID(ctx)
	if !ok {
		return "", fmt.Errorf("cannot get user id from context")
	}

	fileID := uuid.NewRandom().String()

	err := s.gcs.UploadFileToStorage(ctx, request.File, fileID)
	if err != nil {
		return "", fmt.Errorf("cannot load file to object storage: %w", err)
	}

	id, err := s.repo.AddCandidate(userID, request.CandidateName, request.CandidateSurname, fileID)
	if err != nil {
		return "", fmt.Errorf("cannot add candidate to database: %w", err)
	}

	return id, nil
}

// GetRequests returns user requests.
func (s *ReferralService) GetRequests(userID, status string, pageNumber, pageSize int) ([]repository.UserRequests, error) {
	requests, err := s.repo.GetRequests(userID, status, pageNumber, pageSize)
	if err != nil {
		return nil, fmt.Errorf("cannot get user requests: %w", err)
	}

	return requests, nil
}

// DownloadFile downloads file from object storage.
func (s *ReferralService) DownloadFile(ctx context.Context, candidateID string, userID string) error {
	fileID, err := s.repo.GetCVID(candidateID, userID)
	if err != nil {
		return fmt.Errorf("cannot get file id from object storage: %w", err)
	}

	err = s.gcs.DownloadFile(ctx, fileID, fileID)
	if err != nil {
		return fmt.Errorf("cannot download file from object storage: %w", err)
	}

	return nil
}

// UpdateRequest updates request's status.
func (s *ReferralService) UpdateRequest(id, status string) error {
	err := s.repo.UpdateRequest(id, status)
	if errors.Is(err, repository.ErrNoResult) {
		return ErrNoResult
	}
	if err != nil {
		return fmt.Errorf("cannot update user request: %w", err)
	}

	return nil
}
