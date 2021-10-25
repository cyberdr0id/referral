// Package models contains of data models.
package models

// Candidate presents model of a sent candidate.
type Candidate struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Cvosfileid int    `json:"cvosfileid"`
}
