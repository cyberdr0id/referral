package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

var (
	// ErrNoUser handle an error when tyring to get non-database user.
	ErrNoUser = errors.New("user doesn't exists")

	// ErrNoFile handle an error when user try to get non-database CV.
	ErrNoFile = errors.New("there is no file with input id")

	// ErrNoResult presents an error when there are no results for the entered data.
	ErrNoResult = errors.New("there are no results for the entered data")

	// ErrUserAlreadyExists handles an error when user tries to sign up with existing data.
	ErrUserAlreadyExists = errors.New("user already exists")
)

const (
	errorCodeName = "unique_violation"
)

// CreateUser registers a new user.
func (r *Repository) CreateUser(name, password string) (string, error) {
	var id string

	query := `INSERT INTO 
				users(name, password)
			  VALUES
			  	($1, $2)
			  RETURNING 
			  	id;`

	err := r.db.QueryRow(query, name, password).Scan(&id)
	if err, ok := err.(*pq.Error); ok && err.Code.Name() == errorCodeName {
		return "", ErrUserAlreadyExists
	}
	if err != nil {
		return "", fmt.Errorf("cannot add user to database: %w", err)
	}

	return id, nil
}

// GetUser gives user for authorization.
func (r *Repository) GetUser(name string) (User, error) {
	var user User

	query := `SELECT 
				id, name, password, is_admin, created, updated 
			  FROM 
			  	users 
			  WHERE 
			  	name = $1;`

	err := r.db.QueryRow(query, name).Scan(&user.ID, &user.Name, &user.Password, &user.IsAdmin, &user.Created, &user.Updated)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNoUser
	}
	if err != nil {
		return User{}, fmt.Errorf("cannot get user from database: %w", err)
	}

	return user, nil
}
