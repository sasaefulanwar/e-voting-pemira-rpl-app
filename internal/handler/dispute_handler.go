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

func (h *DisputeHandler) SubmitDispute(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(2 << 20)
	if err != nil {
		http.Error(w, "Ukuran file terlalu besar, maks 2MB!", http.StatusBadRequest)
		return
	}

	nim := r.FormValue("nim")
	emailValue := r.Context().Value("email")
	email, ok := emailValue.(string)
	if nim == "" || !ok {
		http.Error(w, "Data tidak lengkap atau unauthorized", http.StatusBadRequest)
		return
	}

	// Ambil file KTM
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File KTM wajib diupload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// ---------------------------------------------------------
	// 🔥 PERBAIKAN NFR-07: MIME TYPE CHECKING (Bukan ekstensi!)
	// ---------------------------------------------------------
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		http.Error(w, "Gagal membaca file", http.StatusInternalServerError)
		return
	}

	// Kembalikan pointer file ke awal setelah dibaca
	file.Seek(0, io.SeekStart)

	fileType := http.DetectContentType(buff)
	if fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/jpg" {
		http.Error(w, "Format file ditolak! Wajib JPG/PNG asli, dilarang memalsukan ekstensi!", http.StatusUnsupportedMediaType)
		return
	}
	// ---------------------------------------------------------

	// Simpan file sementara ke lokal (Nanti kita ganti ke S3/MinIO)
	filename := fmt.Sprintf("%d-%s", time.Now().Unix(), header.Filename)
	filepath := fmt.Sprintf("./uploads/ktm/%s", filename)

	out, err := os.Create(filepath)
	if err != nil {
		http.Error(w, "Gagal menyimpan file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	io.Copy(out, file)

	// Submit sengketa
	err = h.svc.SubmitDispute(nim, email, filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Sengketa berhasil diajukan, akun lu otomatis disuspend sementara!",
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
