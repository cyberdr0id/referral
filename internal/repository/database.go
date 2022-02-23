// Package repository implements types for work with database layer.
package repository

import (
	"database/sql"
	"fmt"
	"github.com/cyberdr0id/referral/internal/config"
	_ "github.com/lib/pq"
)

// NewConnection creates PostgreSQL connection.
func NewConnection(config *config.Database) (*sql.DB, error) {
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		config.Host, config.Port, config.User, config.DatabaseName, config.Password, config.SSLMode)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, fmt.Errorf("cannot create database connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("database created, but cannot be pinged: %w", err)
	}
	return db, nil
}
