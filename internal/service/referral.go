package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/cyberdr0id/referral/internal/repository"
	mycontext "github.com/cyberdr0id/referral/pkg/context"
)

// ErrInvalidParameter presetns an error when user input invalid parameter.
var ErrInvalidParameter = errors.New("invalid parameter")

// ReferralService presents access to referral service via repository.
type ReferralService struct {
	repo *repository.Repository
}

// NewReferralService creates a new instance of ReferralService.
func NewReferralService(repo *repository.Repository) *ReferralService {
	return &ReferralService{
		repo: repo,
	}
}

// SubmitCandidateRequest presents a type for reading data after submitting a candidate.
type SubmitCandidateRequest struct {
	FileName         string
	CandidateName    string
	CandidateSurname string
}

// AddCandidate create request with candidate.
func (s *ReferralService) AddCandidate(ctx context.Context, request SubmitCandidateRequest) (string, error) {
	//TODO: add file to object storage
	fileID := "1"

	userID, ok := mycontext.Get(ctx)
	if !ok {
		return "", fmt.Errorf("cannot get user id from context")
	}

	id, err := s.repo.AddCandidate(userID, request.CandidateName, request.CandidateSurname, fileID)
	if err != nil {
		return "", fmt.Errorf("cannot add candidate to database: %w", err)
	}

	return id, nil
}

// ValidateCandidateSendingRequest validates data after sending a candidate.
func (r *SubmitCandidateRequest) ValidateCandidateSendingRequest() error {
	if len(r.CandidateName) == 0 || len(r.CandidateSurname) == 0 || len(r.FileName) == 0 {
		return fmt.Errorf("%w: wrong length", ErrInvalidParameter)
	}
	fileExp := "([a-zA-Z0-9\\s_\\.\\-\\(\\):])+(.PDF|.pdf)$"
	isRightFile, _ := regexp.MatchString(fileExp, r.FileName)
	if !isRightFile {
		return fmt.Errorf("%w: invalid filename or filetype", ErrInvalidParameter)
	}

	nameSurnameExp := "(^[A-Za-zА-Яа-я]{2,16})?([ ]{0,1})([A-Za-zА-Яа-я]{2,16})?"
	isValid, _ := regexp.MatchString(nameSurnameExp, r.CandidateName)
	if !isValid {
		return fmt.Errorf("%w: name has invalid format", ErrInvalidParameter)
	}

	isValid, _ = regexp.MatchString(nameSurnameExp, r.CandidateSurname)
	if !isValid {
		return fmt.Errorf("%w: surname has invalid format", ErrInvalidParameter)
	}

	return nil
}
