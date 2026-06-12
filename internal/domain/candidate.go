package domain

import "time"

type Candidate struct {
	ID              int64     `json:"id" db:"id"`
	ElectionID      int64     `json:"election_id" db:"election_id"`
	CandidateNumber int       `json:"candidate_number" db:"candidate_number"`
	Name            string    `json:"name" db:"name"`
	Vision          string    `json:"vision" db:"vision"`
	Mission         string    `json:"mission" db:"mission"`
	PhotoURL        string    `json:"photo_url" db:"photo_url"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type Result struct {
	PaslonID int `json:"paslon_id"`
	Votes    int `json:"votes"`
}
