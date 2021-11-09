package repository

import (
	"database/sql"
	"time"
)

// AuthRepository presents methods for user authorization/registration.
type AuthRepository interface {
	CreateUser(name, password string) (string, error)
	GetUser(name, password string) (User, error)
}

// ReferralRepositpry presents methods for candidates manipulating.
type ReferralRepository interface {
	GetRequests(id int) ([]Request, error)
	AddCandidate(name, surname string, fileID int) (string, error)
	GetCVID(id int) (string, error)
	UpdateRequest(id, status string) error
}

// Candidate presents model of a sent candidate.
type Candidate struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	CVOSFileID string `json:"cvosfileid"`
}

// User presents model of user.
type User struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Password string    `json:"password"`
	IsAdmin  bool      `json:"isadmin"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

// Request presents model of request.
type Request struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userid"`
	CandidateID string    `json:"candidateid"`
	Status      string    `json:"status"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

// Repository type which presents connection between database and app logic.
type Repository struct {
	db *sql.DB
	AuthRepository
	ReferralRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
