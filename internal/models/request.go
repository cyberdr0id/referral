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
