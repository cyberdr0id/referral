package repository

import (
	"database/sql"
	"errors"
)

var (
	// ErrNoUser handle an error when tyring to get non-database user.
	ErrNoUser = errors.New("user doesn't exists")
	// ErrNoFile handle an error when user try to get non-database CV.
	ErrNoFile = errors.New("there is no file with input id")
	// ErrNoFile presents an error when there are no results for the entered data.
	ErrNoResult = errors.New("there are no results for the entered data")
)

// CreateUser registers a new user.
func (r *Repository) CreateUser(name, password string) (string, error) {
	var id string

	row := r.DB.QueryRow("INSERT INTO USERS(NAME, PASSWORD) VALUES($1, $2) RETURNING ID;", name, password)
	err := row.Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

// GetUser gives user for authorization.
func (r *Repository) GetUser(name string) (User, error) {
	var user User

	row := r.DB.QueryRow("SELECT ID, NAME, PASSWORD, ISADMIN, CREATED, UPDATED FROM USERS WHERE NAME=$1;", name)
	err := row.Scan(&user.ID, &user.Name, &user.Password, &user.IsAdmin, &user.Created, &user.Updated)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNoUser
	}
	if err != nil {
		return User{}, err
	}

	return user, nil
}
