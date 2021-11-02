package repository

// GetRequest gives user requests by id
func (r *Repository) GetRequests(id int) ([]Request, error) {
	var requests []Request

	rows, err := r.db.Query("SELECT ID, USERID, CANDIDATEID, STATUS, CREATED, UPDATED FROM REQUESTS WHERE USERID = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		request := Request{}
		err := rows.Scan(&request.ID, &request.UserID, &request.CandidateID, &request.Status, &request.Created, &request.Updated)
		if err != nil {
			return nil, err
		}

		requests = append(requests, request)
	}
	return requests, nil
}

// TODO: change return value to Candidate object if only ID is bad and change response in swagger file
// AddCandidate adds submitted candidate.
func (r *Repository) AddCandidate(name, surname string, fileID int) (int, error) {
	var requestID int

	row := r.db.QueryRow("INSERT INTO CANDIDATES(NAME, SURNAME, CVOSFILEID) VALUES($1, $2, $3)", name, surname, fileID)
	if err := row.Scan(&requestID); err != nil {
		return 0, err
	}

	return requestID, nil
}

// GetCVID returns cv file id from object storage.
func (r *Repository) GetCVID(id int) (int, error) {
	var fileID int

	err := r.db.QueryRow("SELECT CVOSFILEID FROM CANDIDATES WHERE ID = $1", id).Scan(&fileID)
	if err != nil {
		return 0, err
	}

	return fileID, nil
}
