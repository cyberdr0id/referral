// Package repository implements types for work with database layer.
package repository

import (
	"database/sql"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

// NewConnection creates PostgreSQL connection.
func NewConnection() (*sql.DB, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to load database config: %w", err)
	}

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

func loadConfig() (*repositoryConfig, error) {
	var c repositoryConfig

	err := envconfig.Process("db", &c)
	if err != nil {
		return nil, fmt.Errorf("unable to read database config: %w (missin env variables?)", err)
	}

	return &c, nil
}
