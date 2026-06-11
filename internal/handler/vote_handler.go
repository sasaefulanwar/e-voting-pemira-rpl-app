package handler

import (
	"encoding/json"
	"net/http"
	"pemira-rpl/internal/service" // Sesuaikan path module lu
)

type VoteHandler struct {
	voteService service.VoteService
}

func NewVoteHandler(vs service.VoteService) *VoteHandler {
	return &VoteHandler{voteService: vs}
}

func (h *VoteHandler) CastVote(w http.ResponseWriter, r *http.Request) {

	voterID, ok := r.Context().Value("voter_id").(int)
	if !ok {
		http.Error(w, `{"error": "Lu siapa cuy? Login dulu!"}`, http.StatusUnauthorized)
		return
	}

	var req struct {
		CandidateID int `json:"candidate_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Format data salah!"}`, http.StatusBadRequest)
		return
	}

	err := h.voteService.CastVote(voterID, req.CandidateID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Voting berhasil! Suara lu udah masuk database cuy."}`))
}
