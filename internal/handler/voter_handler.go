package handler

import (
	"encoding/json"
	"net/http"
	"pemira-rpl/internal/dto"
	"pemira-rpl/internal/service"
	"pemira-rpl/internal/utils"
)

type VoterHandler struct {
	srv service.VoterService
}

func NewVoterHandler(srv service.VoterService) *VoterHandler {
	return &VoterHandler{srv: srv}
}

func (h *VoterHandler) BindNIM(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method tidak diizinkan"})
		return
	}

	var req dto.BindNIMRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Format data rusak"})
		return
	}

	realEmail, ok := r.Context().Value("email").(string)
	if !ok || realEmail == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Gagal mendapatkan email dari sesi"})
		return
	}

	res, err := h.srv.ProcessBinding(req, realEmail)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	jwtToken, err := utils.GenerateJWT(
		realEmail,
		req.Nama,
		req.NIM,
	)

	if err != nil {
		w.WriteHeader(
			http.StatusInternalServerError,
		)

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "gagal membuat JWT baru",
			},
		)

		return
	}

	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "jwt_token",
			Value:    jwtToken,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
		},
	)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

}
