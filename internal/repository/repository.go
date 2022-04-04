package repository

import (
	"database/sql"
)

// AuthRepository presents methods for user authorization/registration.
type AuthRepository interface {
	CreateUser(name, password string) (string, error)
	GetUser(name, password string) (User, error)
}

// ReferralRepository presents methods for candidates manipulating.
type ReferralRepository interface {
	GetRequests(id int) ([]UserRequests, error)
	AddCandidate(name, surname string, fileID int) (string, error)
	GetCVID(id int) (string, error)
	UpdateRequest(id, status string) error
}

// Repository type which presents connection between database and app logic.
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new instance of Repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
