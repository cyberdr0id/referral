// Package repository implements types for work with database layer.
package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// DatabaseConfig represents a type that contains database configuration data.
type DatabaseConfig struct {
	Host         string
	User         string
	Password     string
	DatabaseName string
	Port         string
	SSLMode      string
}

// NewConnection creates PostgreSQL connection.
func NewConnection(config DatabaseConfig) (*sql.DB, error) {
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
