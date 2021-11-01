package repository

import (
	"database/sql"
	"time"
)

// Repository type which presents connection between database and app logic.
type Repository struct {
	db *sql.DB
}

// Candidate presents model of a sent candidate.
type Candidate struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Cvosfileid int    `json:"cvosfileid"`
}

// User presents model of user.
type User struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Password string    `json:"password"`
	IsAdmin  bool      `json:"isadmin"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
