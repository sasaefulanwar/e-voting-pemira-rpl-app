package domain

import "time"

type Candidate struct {
	ID              int64 `json:"id"`
	ElectionID      int64 `json:"election_id"`
	CandidateNumber int   `json:"candidate_number"`

	ChairmanName     string `json:"chairman_name"`
	ViceChairmanName string `json:"vice_chairman_name"`

	Vision  string `json:"vision"`
	Mission string `json:"mission"`

	PhotoURL string `json:"photo_url"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Result struct {
	PaslonID int `json:"paslon_id"`
	Votes    int `json:"votes"`
}
