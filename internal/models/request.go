package models

import "time"

// Request presents model of request.
type Request struct {
	ID          int       `json:"id"`
	UserID      int       `json:"userid"`
	CandidateID int       `json:"candidateid"`
	Status      string    `json:"status"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

// AuthRequest presents request for authorization.
type AuthRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// CandidateSendingRequest presents request for sending candidate.
type CandidateSendingRequest struct {
	FileName         string
	CandidateName    string
	CandidateSurname string
}
