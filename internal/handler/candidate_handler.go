package handler

import (
	"encoding/json"
	"net/http"
	"pemira-rpl/internal/service"
)

type CandidateHandler struct {
	svc service.CandidateService
}

func NewCandidateHandler(svc service.CandidateService) *CandidateHandler {
	return &CandidateHandler{svc: svc}
}

func (h *CandidateHandler) GetAll(
	w http.ResponseWriter,
	r *http.Request,
) {

	candidates, err := h.svc.GetAll()

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(candidates)
}

func (h *CandidateHandler) GetResults(
	w http.ResponseWriter,
	r *http.Request,
) {

	email, _ := r.Context().
		Value("email").(string)

	isAdmin :=
		email ==
			"himarpl@polindra.ac.id"

	results, err :=
		h.svc.GetResults(
			isAdmin,
		)

	if err != nil {

		w.WriteHeader(
			http.StatusForbidden,
		)

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": err.Error(),
			},
		)

		return
	}

	json.NewEncoder(w).Encode(results)
}
