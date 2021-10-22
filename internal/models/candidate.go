package models

type Candidate struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Cvosfileid int    `json:"cvosfileid"`
}
