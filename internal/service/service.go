package service

import (
	"context"
	"errors"

	"github.com/cyberdr0id/referral/internal/repository"
	"github.com/dgrijalva/jwt-go"
)

var (
	// ErrNoUser handle an error when tyring to get non-database user.
	ErrNoUser = errors.New("user doesn't exists")
	// ErrNoFile handle an error when user try to get non-database CV.
	ErrNoFile = errors.New("there is no file with input id")
	// ErrNoResult presents an error when there are no results for the entered data.
	ErrNoResult = errors.New("there are no results for the entered data")
	// ErrUserAlreadyExists handles an error when user tries to sign up with existing data.
	ErrUserAlreadyExists = errors.New("user already exists")
)

// Auth presents interface for authorization.
type Auth interface {
	LogIn(name, password string) (string, error)
	SignUp(name, password string) (string, error)
	ParseToken(_token string) (*jwt.StandardClaims, error)
}

// Referral presents a type of CV interaction.
type Referral interface {
	GetRequests(id, t string) ([]repository.Request, error)
	AddCandidate(ctx context.Context, request SubmitCandidateRequest) (string, error)
	DownloadFile(id string) (string, error)
	UpdateRequest(id, status string) error
}
