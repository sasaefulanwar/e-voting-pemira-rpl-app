package handler

import (
	"encoding/json"
	"net/http"
	"pemira-rpl/internal/service"
)

type ElectionHandler struct {
	srv service.ElectionService
}

func NewElectionHandler(srv service.ElectionService) *ElectionHandler {
	return &ElectionHandler{srv: srv}
}

// GET /api/v1/election/status
func (h *ElectionHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	stats, err := h.srv.GetStatistics()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "gagal mengambil status election",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": stats.ElectionStatus, // open / closed
	})
}

// POST /api/v1/admin/election/open
func (h *ElectionHandler) OpenElection(w http.ResponseWriter, r *http.Request) {
	err := h.srv.OpenElection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "gagal membuka election",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Election berhasil dibuka",
	})
}

func (h *ElectionHandler) CloseElection(w http.ResponseWriter, r *http.Request) {
	err := h.srv.CloseElection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "gagal menutup election",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Election berhasil ditutup",
	})
}
