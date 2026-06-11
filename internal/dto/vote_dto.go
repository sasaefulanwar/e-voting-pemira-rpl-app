package dto

type VoteRequest struct {
	ElectionID int `json:"election_id"`
	PaslonID   int `json:"paslon_id"`
}
