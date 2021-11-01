// Package repository implements types for work with database layer.
package repository

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

var (
	// errSqlConnection presents an error when sql connection cannot be created.
	errSQLConnection = errors.New("cannot create sql connection")

	// errPingDatabase presents an error when connection is created but database cannot be pinged.
	errPingDatabase = errors.New("connection is created but database cannot be pinged")
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
		return nil, errSQLConnection
	}

	if err = db.Ping(); err != nil {
		return nil, errPingDatabase
	}
	return db, nil
}
