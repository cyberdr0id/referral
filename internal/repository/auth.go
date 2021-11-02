package repository

// CreateUser registers a new user.
func (r *Repository) CreateUser(name, password string) (int, error) {
	var id int

	row := r.db.QueryRow("INSERT INTO USERS(NAME, PASSWORD) VALUES($1, $2)", name, password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// GetUser gives user for authorization.
func (r *Repository) GetUser(name, password string) (User, error) {
	var user User

	row := r.db.QueryRow("SELECT ID, NAME, PASSWORD, ISADMIN, CREATED, UPDATED FROM USERS WHERE NAME=$1 AND PASSWORD=$1", name, password)
	if err := row.Scan(&user.ID, &user.Name, &user.Password, &user.IsAdmin, &user.Created, &user.Updated); err != nil {
		return User{}, err
	}

	return user, nil
}
