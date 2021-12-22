package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	mycontext "github.com/cyberdr0id/referral/internal/context"
	"github.com/cyberdr0id/referral/internal/repository"
	"github.com/cyberdr0id/referral/internal/storage"
)

// ErrInvalidParameter presetns an error when user input invalid parameter.
var ErrInvalidParameter = errors.New("invalid parameter")

// ReferralService presents access to referral service via repository.
type ReferralService struct {
	repo *repository.Repository
	s3   *storage.Storage
}

// NewReferralService creates a new instance of ReferralService.
func NewReferralService(repo *repository.Repository, s3 *storage.Storage) *ReferralService {
	return &ReferralService{
		repo: repo,
		s3:   s3,
	}
}

// SubmitCandidateRequest presents a type for reading data after submitting a candidate.
type SubmitCandidateRequest struct {
	File             multipart.File
	CandidateName    string
	CandidateSurname string
}

// AddCandidate create request with candidate.
func (s *ReferralService) AddCandidate(ctx context.Context, request SubmitCandidateRequest) (string, error) {
	userID, ok := mycontext.GetUserID(ctx)
	if !ok {
		return "", fmt.Errorf("cannot get user id from context")
	}

	fileID := "1"

	err := s.s3.UploadFileToStorage(request.File, fileID)
	if err != nil {
		return "", fmt.Errorf("cannot load file to object storage: %w", err)
	}

	id, err := s.repo.AddCandidate(userID, request.CandidateName, request.CandidateSurname, fileID)
	if err != nil {
		return "", fmt.Errorf("cannot add candidate to database: %w", err)
	}

	return id, nil
}
