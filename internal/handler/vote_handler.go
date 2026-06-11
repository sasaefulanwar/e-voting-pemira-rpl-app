package handler

import (
	"encoding/json"
	"net/http"
	"pemira-rpl/internal/service"
)

type VoteHandler struct {
	voteService service.VoteService
}

func NewVoteHandler(vs service.VoteService) *VoteHandler {
	return &VoteHandler{
		voteService: vs,
	}
}

type VoteRequest struct {
	ElectionID int `json:"election_id"`
	PaslonID   int `json:"paslon_id"`
}

func (h *VoteHandler) CastVote(
	w http.ResponseWriter,
	r *http.Request,
) {

	nim, ok := r.Context().
		Value("nim").(string)

	if !ok || nim == "" {
		http.Error(
			w,
			`{"error":"NIM tidak ditemukan di token"}`,
			http.StatusUnauthorized,
		)
		return
	}

	var req VoteRequest

	if err := json.NewDecoder(r.Body).
		Decode(&req); err != nil {

		http.Error(
			w,
			`{"error":"request tidak valid"}`,
			http.StatusBadRequest,
		)
		return
	}

	err := h.voteService.CastVote(
		nim,
		req.ElectionID,
		req.PaslonID,
	)

	if err != nil {

		w.Header().
			Set(
				"Content-Type",
				"application/json",
			)

		w.WriteHeader(
			http.StatusConflict,
		)

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": err.Error(),
			},
		)

		return
	}

	json.NewEncoder(w).Encode(
		map[string]string{
			"message": "vote berhasil",
		},
	)
}
