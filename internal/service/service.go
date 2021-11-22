package service

import (
	"github.com/cyberdr0id/referral/internal/repository"
)

type Auth interface {
	SignUp(string, string) (string, error)
	LogIn(string, string) (string, error)
}

type Referral interface {
	SendCandidate(string, string, string, string) (string, error)
	GetRequests(string, string) ([]repository.Request, error)
	DownloadCV(string) (string, error)
	UpdateRequest(string, string, string) error
}
