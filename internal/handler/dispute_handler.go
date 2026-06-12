package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"pemira-rpl/internal/domain"
)

type DisputeService interface {
	SubmitDispute(
		nim,
		email,
		ktmPath string,
	) error

	GetAllDisputes() (
		[]domain.Dispute,
		error,
	)

	ApproveDispute(
		id int64,
	) error

	RejectDispute(
		id int64,
	) error
}

type DisputeHandler struct {
	svc DisputeService
}

func NewDisputeHandler(
	svc DisputeService,
) *DisputeHandler {

	return &DisputeHandler{
		svc: svc,
	}
}

func (
	h *DisputeHandler,
) GetAllDisputes(
	w http.ResponseWriter,
	r *http.Request,
) {

	disputes,
		err :=
		h.svc.GetAllDisputes()

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

	json.NewEncoder(
		w,
	).Encode(
		disputes,
	)
}

func (h *DisputeHandler) SubmitDispute(
	w http.ResponseWriter,
	r *http.Request,
) {
	// Maksimal file 5MB
	err := r.ParseMultipartForm(5 << 20)
	if err != nil {
		http.Error(w, "Gagal parsing form data", http.StatusBadRequest)
		return
	}

	nim := r.FormValue("nim")
	if nim == "" {
		http.Error(w, "NIM wajib diisi", http.StatusBadRequest)
		return
	}

	emailValue := r.Context().Value("email")
	if emailValue == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	email, ok := emailValue.(string)
	if !ok {
		http.Error(w, "invalid email context", http.StatusUnauthorized)
		return
	}

	// Ambil file KTM
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File KTM wajib diupload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Simpan file ke folder lokal ./uploads/ktm
	filename := fmt.Sprintf("%d-%s", time.Now().Unix(), header.Filename)
	filepath := fmt.Sprintf("./uploads/ktm/%s", filename)

	out, err := os.Create(filepath)
	if err != nil {
		http.Error(w, "Gagal menyimpan file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Gagal menyimpan file", http.StatusInternalServerError)
		return
	}

	// Submit sengketa
	err = h.svc.SubmitDispute(nim, email, filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Sengketa berhasil diajukan",
		"ktm":     filename,
	})
}

func (
	h *DisputeHandler,
) GetPending(
	w http.ResponseWriter,
	r *http.Request,
) {

	disputes, err :=
		h.svc.GetAllDisputes()

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	json.NewEncoder(
		w,
	).Encode(
		disputes,
	)
}

func (
	h *DisputeHandler,
) ApproveDispute(
	w http.ResponseWriter,
	r *http.Request,
) {

	idStr :=
		r.URL.Query().
			Get("id")

	id,
		err :=
		strconv.ParseInt(
			idStr,
			10,
			64,
		)

	if err != nil {

		http.Error(
			w,
			"invalid id",
			http.StatusBadRequest,
		)

		return
	}

	err =
		h.svc.
			ApproveDispute(id)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	json.NewEncoder(
		w,
	).Encode(
		map[string]string{
			"message": "dispute approved",
		},
	)
}

func (
	h *DisputeHandler,
) RejectDispute(
	w http.ResponseWriter,
	r *http.Request,
) {

	idStr :=
		r.URL.Query().
			Get("id")

	id,
		err :=
		strconv.ParseInt(
			idStr,
			10,
			64,
		)

	if err != nil {

		http.Error(
			w,
			"invalid id",
			http.StatusBadRequest,
		)

		return
	}

	err =
		h.svc.
			RejectDispute(id)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	json.NewEncoder(
		w,
	).Encode(
		map[string]string{
			"message": "dispute rejected",
		},
	)
}
