package repository

import (
	"database/sql"
	"errors"
)

// GetRequests gives user requests by id.
func (r *Repository) GetRequests(id string, t string) ([]Request, error) {
	var requests []Request

	if t == "" {
		t = "ID"
	}

	query := `SELECT id, userid, candidateid, status, created, updated 
			  FROM requests WHERE userid = $1 
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

// TODO: change return value to Candidate object if only ID is bad and change response in swagger file
// AddCandidate adds submitted candidate.
func (r *Repository) AddCandidate(name, surname, fileID string) (string, error) {
	var requestID string

	query := `INSERT INTO candidates(name, surname, cvosfileid) 
			  VALUES($1, $2, $3) RETURNING id;`

	err := r.db.QueryRow(query, name, surname, fileID).Scan(&requestID)
	if err != nil {
		return "", err
	}

	return requestID, nil
}

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

	query := `SELECT cvosfileid
			  FROM candidates
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

func (r *Repository) CreateRequest(userID, candidateID string) (string, error) {
	var id string

	query := `INSERT INTO 
			  	requests(userid, candidateid) 
			  VALUES($1, $2) RETURNING id;`

	row := r.db.QueryRow(query, userID, candidateID)
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (r *Repository) IsUserRequest(userID, requestID string) error {
	var id string

	query := `SELECT id
			  FROM requests 
			  WHERE userid = $1 AND id = $2`

	err := r.db.QueryRow(query, userID, requestID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNoAccess
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) IsUserAdmin(userID string) (bool, error) {
	var isadmin string

	query := `SELECT isadmin
			  FROM users 
			  WHERE id = $1`

	err := r.db.QueryRow(query, userID).Scan(&isadmin)
	if errors.Is(err, sql.ErrNoRows) {
		return false, ErrNoAccess
	}
	if err != nil {
		return false, err
	}

	return isadmin == "true", nil
}
