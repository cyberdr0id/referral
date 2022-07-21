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
	Author  author `json:"author"`
}

type author struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"isadmin"`
}

// GetRequests gives user requests by id.
func (r *Repository) GetRequests(id, status string, pageNumber, pageSize int) ([]UserRequests, error) {
	var requests []UserRequests
	var whereVal []interface{}

	query := `
			SELECT
				id, candidate_name, candidate_surname, status, updated
			FROM
				requests%s
			LIMIT $%d
			OFFSET $%d
			`

	offset := (pageNumber - 1) * pageSize

	if id == "" {
		if status == "" {
			query = fmt.Sprintf(query, "", 1, 2)
			whereVal = append(whereVal, pageSize, offset)
		} else {
			query = fmt.Sprintf(query, " WHERE status = $1 ", 2, 3)
			whereVal = append(whereVal, status, pageSize, offset)
		}
	} else {
		if status == "" {
			query = fmt.Sprintf(query, " WHERE author_id = $1 ", 2, 3)
			whereVal = append(whereVal, id, pageSize, offset)
		} else {
			query = fmt.Sprintf(query, " WHERE author_id = $1 AND status = $2 ", 3, 4)
			whereVal = append(whereVal, id, status, pageSize, offset)
		}
	}

	rows, err := r.db.Query(query, whereVal...)
	if err != nil {
		return nil, fmt.Errorf("error with query executing: %w", err)
	}

	// user, _ := r.GetUser(id)

	for rows.Next() {
		request := UserRequests{}

		if err := rows.Scan(
			&request.ID,
			&request.Name,
			&request.Surname,
			&request.Status,
			&request.Updated,
		); err != nil {
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

	query := `INSERT INTO 
				requests(author_id, candidate_name, candidate_surname, cv_file_id) 
			  VALUES
			  	($1, $2, $3, $4)
			  RETURNING id;`

	err := r.db.QueryRow(query, userID, name, surname, fileID).Scan(&requestID)
	if err != nil {
		return "", fmt.Errorf("cannot add candidate to database: %w", err)
	}

	return requestID, nil
}

// UpdateRequest updates user request status.
func (r *Repository) UpdateRequest(id, newState string) error {
	query := `UPDATE 
				requests 
			  SET 
			  	status = $1
			  WHERE 
			  	id = $2;`

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
func (r *Repository) GetCVID(candidateID, userID string) (string, error) {
	var fileID string

	query := `SELECT
				cv_file_id
			  FROM 
			  	requests
			  WHERE 
			  	id = $1`

	whereVal := []interface{}{
		candidateID,
	}

	if userID != "" {
		query = fmt.Sprintf("%s AND author_id = $2", query)
		whereVal = append(whereVal, userID)
	}

	err := r.db.QueryRow(query, whereVal...).Scan(&fileID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNoFile
	}
	if err != nil {
		return "", fmt.Errorf("cannot get cv id from database: %w", err)
	}

	return fileID, nil
}
