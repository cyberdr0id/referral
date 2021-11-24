package service

import (
	"errors"

	"github.com/cyberdr0id/referral/internal/repository"
)

// ErrInvalidParameter presetns an error when user input invalid parameter.
var ErrInvalidParameter = errors.New("invalid parameter")

type Auth interface {
	LogIn(name, password string) (string, string, error)
	SignUp(name, password string) (string, error)
}

type Referral interface {
	GetRequests(id, t string) ([]repository.Request, error)
	AddCandidate(name, surname, fileID string) (string, error)
	GetCVID(id string) (string, error)
	UpdateRequest(id, status string) error
}
