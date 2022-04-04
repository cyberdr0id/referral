package service

import (
	"context"

	"github.com/cyberdr0id/referral/internal/repository"
	myjwt "github.com/cyberdr0id/referral/pkg/jwt"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	// ErrNoUser handle an error when trying to get non-database user.
	ErrNoUser = Error("user doesn't exists")

	// ErrNoFile handle an error when user try to get non-database CV.
	ErrNoFile = Error("there is no file with input id")

	// ErrNoResult presents an error when there are no results for the entered data.
	ErrNoResult = Error("there are no results for the entered data")

	// ErrUserAlreadyExists handles an error when user tries to sign up with existing data.
	ErrUserAlreadyExists = Error("user already exists")
)

// Auth presents interface for authorization and registration actions.
type Auth interface {
	LogIn(name, password string) (string, error)
	SignUp(name, password string) (string, error)
	ParseToken(token string) (*myjwt.Claims, error)
}

// Referral presents a type of CV interaction.
type Referral interface {
	GetRequests(userID, status string, pageNumber, pageSize int) ([]repository.UserRequests, error)
	AddCandidate(ctx context.Context, request SubmitCandidateRequest) (string, error)
	DownloadFile(ctx context.Context, id string, userID string) (string, error)
	UpdateRequest(id, status string) error
}
