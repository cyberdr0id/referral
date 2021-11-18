package service

import "github.com/cyberdr0id/referral/internal/repository"

type Auth interface {
	LogIn(name, password string) (string, error)
	SignUp(name, password string) (string, error)
}

type Referral interface {
	GetRequests(id int) ([]repository.Request, error)
	AddCandidate(name, surname string, fileID int) (string, error)
	GetCVID(id int) (string, error)
	UpdateRequest(id, status string) error
}
