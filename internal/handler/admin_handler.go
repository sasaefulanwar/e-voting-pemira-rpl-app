package handler

import (
	"encoding/json"
	"net/http"

	"pemira-rpl/internal/service"
)

type AdminHandler struct {
	svc service.ElectionService
}

func NewAdminHandler(
	svc service.ElectionService,
) *AdminHandler {
	return &AdminHandler{
		svc: svc,
	}
}

func (h *AdminHandler) OpenElection(
	w http.ResponseWriter,
	r *http.Request,
) {

	if err := h.svc.OpenElection(); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	json.NewEncoder(w).Encode(
		map[string]string{
			"message": "pemilu dibuka",
		},
	)
}

func (h *AdminHandler) CloseElection(
	w http.ResponseWriter,
	r *http.Request,
) {

	if err := h.svc.CloseElection(); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	json.NewEncoder(w).Encode(
		map[string]string{
			"message": "pemilu ditutup",
		},
	)
}

func (h *AdminHandler) GetStatistics(
	w http.ResponseWriter,
	r *http.Request,
) {

	stats,
		err := h.svc.GetStatistics()

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	json.NewEncoder(w).Encode(stats)
}
