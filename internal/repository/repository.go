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
	ID         string
	Name       string
	Surname    string
	CVOSFileID string
}

// User presents model of user.
type User struct {
	ID       string
	Name     string
	Password string
	IsAdmin  bool
	Created  time.Time
	Updated  time.Time
}

// Request presents model of request.
type Request struct {
	ID          string
	UserID      string
	CandidateID string
	Status      string
	Created     time.Time
	Updated     time.Time
}

// Repository type which presents connection between database and app logic.
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
