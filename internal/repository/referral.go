package repository

import (
	"database/sql"
	"errors"
)

// GetRequests gives user requests by id.
func (r *Repository) GetRequests(id string, t string) ([]Request, error) {
	var requests []Request

	if t == "" {
		t = "Updated"
	}

	query := `SELECT id, user_id, candidate_id, status, created, updated 
			  FROM requests WHERE user_id = $1 
			  ORDER BY $2`

	rows, err := r.db.Query(query, id, t)
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		request := Request{}

		if err := rows.Scan(&request.ID, &request.UserID, &request.CandidateID,
			&request.Status, &request.Created, &request.Updated); err != nil {
			return nil, err
		}

		requests = append(requests, request)
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
		return "", err
	}

	return requestID, nil
}

// UpdateRequest updates user request status.
func (r *Repository) UpdateRequest(id, newState string) error {
	query := `UPDATE requests 
			  SET status = $1
			  WHERE id = $2;`

	rows, err := r.db.Exec(query, newState, id)
	if err != nil {
		return err
	}

	n, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNoResult
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
		return "", err
	}

	return fileID, nil
}
