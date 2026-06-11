package domain

import "time"

type Ballot struct {
	ID         int       `json:"id"`
	ElectionID int       `json:"election_id"`
	HashedNIM  string    `json:"hashed_nim"`
	PaslonID   int       `json:"paslon_id"`
	CreatedAt  time.Time `json:"created_at"`
}
