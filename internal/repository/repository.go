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

// ReferralRepository presents methods for candidates manipulating.
type ReferralRepository interface {
	GetRequests(id int) ([]Request, error)
	AddCandidate(name, surname string, fileID int) (string, error)
	GetCVID(id int) (string, error)
	UpdateRequest(id, status string) error
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

// Config
type repositoryConfig struct {
	Host         string `envconfig:"DB_HOST"`
	User         string `envconfig:"DB_USER"`
	Password     string `envconfig:"DB_PASSWORD"`
	DatabaseName string `envconfig:"DB_NAME"`
	Port         string `envconfig:"DB_PORT"`
	SSLMode      string `envconfig:"DB_SSLMODE"`
}

// Repository type which presents connection between database and app logic.
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new instance of Repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
