package repository

import (
	"database/sql"
	"errors"
	"fmt"
)

// GetRequests gives user requests by id.
func (r *Repository) GetRequests(id string, t string) ([]Request, error) {
	var requests []Request

	if t == "" {
		t = "ID"
	}

	rows, err := r.DB.Query("SELECT ID, USERID, CANDIDATEID, STATUS, CREATED, UPDATED FROM REQUESTS WHERE USERID = $1 ORDER BY $2", id, t)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		request := Request{}

		if err := rows.Scan(&request.ID, &request.UserID, &request.CandidateID, &request.Status, &request.Created, &request.Updated); err != nil {
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

	row := r.DB.QueryRow("INSERT INTO CANDIDATES(NAME, SURNAME, CVOSFILEID) VALUES($1, $2, $3) RETURNING ID;", name, surname, fileID)
	if err := row.Scan(&requestID); err != nil {
		return "", err
	}

	return requestID, nil
}

func (r *Repository) UpdateRequest(id, newState string) error {
	fmt.Printf("UPDATE REQUESTS SET STATUS = %s WHERE ID = %s \n", newState, id)

	rows, err := r.DB.Exec("UPDATE REQUESTS SET STATUS = $1 WHERE ID = $2;", newState, id)
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

	err := r.DB.QueryRow("SELECT CVOSFILEID FROM CANDIDATES WHERE ID = $1", id).Scan(&fileID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNoFile
	}
	if err != nil {
		return "", err
	}

	return fileID, nil
}
