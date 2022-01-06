package repository

import (
	"database/sql"
	"errors"
	"fmt"
)

// UserRequests presents a type for user requests data.
type UserRequests struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Status  string `json:"status"`
	Updated string `json:"updated"`
}

// GetRequests gives user requests by id.
func (r *Repository) GetRequests(id string, status string) ([]UserRequests, error) {
	var requests []UserRequests
	var query string

	user, _ := r.GetUser(id)

	if user.IsAdmin {
		query = "SELECT id, candidate_name, candidate_surname, status, updated FROM requests"
	} else {
		if status == "" {
			query = fmt.Sprintf(`SELECT id, candidate_name, candidate_surname, status, updated 
					 FROM requests WHERE author_id = %s`, id)
		} else {
			query = fmt.Sprintf(`SELECT id, candidate_name, candidate_surname, status, updated 
					 FROM requests WHERE author_id = %s AND requests.status = '%s'`, id, status)
		}
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error with query executing: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		request := UserRequests{}

		if err := rows.Scan(&request.ID, &request.Name, &request.Surname,
			&request.Status, &request.Updated); err != nil {
			return nil, fmt.Errorf("cannot get requests information: %w", err)
		}

		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error with result set: %w", err)
	}

	return requests, nil
}

// AddCandidate adds submitted candidate.
func (r *Repository) AddCandidate(userID, name, surname, fileID string) (string, error) {
	var requestID string

	query := `INSERT INTO requests(author_id, candidate_name, candidate_surname, cv_file_id) 
			  VALUES($1, $2, $3, $4) RETURNING id;`

	err := r.db.QueryRow(query, userID, name, surname, fileID).Scan(&requestID)
	if err != nil {
		return "", fmt.Errorf("cannot add candidate to database: %w", err)
	}

	return requestID, nil
}

// UpdateRequest updates user request status.
func (r *Repository) UpdateRequest(id, newState string) error {
	query := `UPDATE requests 
			  SET status = $1
			  WHERE id = $2;`

	rows, err := r.db.Exec(query, newState, id)

	n, _ := rows.RowsAffected()
	if n == 0 {
		return ErrNoResult
	}
	if err != nil {
		return fmt.Errorf("cannot update user request: %w", err)
	}

	return nil
}

// GetCVID returns cv file id from object storage.
func (r *Repository) GetCVID(id string) (string, error) {
	var fileID string

	query := `SELECT cv_file_id
			  FROM requests
			  WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(&fileID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNoFile
	}
	if err != nil {
		return "", fmt.Errorf("cannot get cv id from database: %w", err)
	}

	return fileID, nil
}
