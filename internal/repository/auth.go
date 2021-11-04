package repository

import (
	"fmt"
)

// CreateUser registers a new user.
func (r *Repository) CreateUser(name, password string) (string, error) {
	var id string

	row := r.db.QueryRow("INSERT INTO USERS(NAME, PASSWORD) VALUES($1, $2) RETURNING ID;", name, password)
	err := row.Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

// GetUser gives user for authorization.
func (r *Repository) GetUser(name string) (User, error) {
	var user User
	fmt.Println(name)
	row := r.db.QueryRow("SELECT ID, NAME, PASSWORD, ISADMIN, CREATED, UPDATED FROM USERS WHERE NAME=$1;", name)
	err := row.Scan(&user.ID, &user.Name, &user.Password, &user.IsAdmin, &user.Created, &user.Updated)
	if err != nil {
		return User{}, err
	}
	fmt.Println(user)
	return user, nil
}
