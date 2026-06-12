package dto

type AdminStatisticsResponse struct {
	TotalVoters    int    `json:"total_voters"`
	Voted          int    `json:"voted"`
	NotVoted       int    `json:"not_voted"`
	ElectionStatus string `json:"election_status"`
}
