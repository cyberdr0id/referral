package service

import (
	"errors"

	"github.com/cyberdr0id/referral/internal/repository"
)

// ErrInvalidParameter presetns an error when user input invalid parameter.
var ErrInvalidParameter = errors.New("invalid parameter")

var (
	// ErrNoUser handle an error when tyring to get non-database user.
	ErrNoUser = errors.New("user doesn't exists")
	// ErrNoFile handle an error when user try to get non-database CV.
	ErrNoFile = errors.New("there is no file with input id")
	// ErrNoFile presents an error when there are no results for the entered data.
	ErrNoResult = errors.New("there are no results for the entered data")
	// ErrUserAlreadyExists handles an error when user tries to sign up with existing data.
	ErrUserAlreadyExists = errors.New("user already exists")
)

type Auth interface {
	LogIn(name, password string) (string, string, error)
	CreateUser(name, password string) (string, error)
}

type Referral interface {
	GetRequests(id, t string) ([]repository.Request, error)
	AddCandidate(name, surname, fileID string) (string, error)
	GetCVID(id string) (string, error)
	UpdateRequest(id, status string) error
}
