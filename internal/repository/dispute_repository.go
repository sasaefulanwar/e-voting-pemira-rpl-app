package repository

import (
	"database/sql"
	"errors"
	"os"
	"pemira-rpl/internal/domain"
	"pemira-rpl/internal/utils"
)

type DisputeRepository interface {
	Submit(dispute domain.Dispute) error
	GetByID(id int64) (*domain.Dispute, error)
	Approve(id int64) error
	Reject(id int64) error
	GetAll() ([]domain.Dispute, error)
	ApproveAndResolveTransaction(disputeID int64) error
}

type disputeRepository struct {
	db *sql.DB
}

func NewDisputeRepository(db *sql.DB) DisputeRepository {
	return &disputeRepository{db: db}
}

func (r *disputeRepository) Submit(dispute domain.Dispute) error {
	_, err := r.db.Exec(`
		INSERT INTO sengketa_nim(nim_sengketa, email_pelapor, path_foto_ktm, status_proses)
		VALUES($1, $2, $3, 'pending')
	`, dispute.NIM, dispute.ReporterEmail, dispute.KTMPath)
	return err
}

func (r *disputeRepository) GetByID(id int64) (*domain.Dispute, error) {
	var d domain.Dispute
	err := r.db.QueryRow(`
		SELECT id, nim_sengketa, email_pelapor, path_foto_ktm, status_proses
		FROM sengketa_nim
		WHERE id=$1
	`, id).Scan(&d.ID, &d.NIM, &d.ReporterEmail, &d.KTMPath, &d.Status)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *disputeRepository) Approve(id int64) error {
	_, err := r.db.Exec(`
		UPDATE sengketa_nim
		SET status_proses='approved'
		WHERE id=$1
	`, id)
	return err
}

func (r *disputeRepository) Reject(id int64) error {
	_, err := r.db.Exec(`
		UPDATE sengketa_nim
		SET status_proses='rejected'
		WHERE id=$1
	`, id)
	return err
}

func (r *disputeRepository) GetAll() ([]domain.Dispute, error) {
	rows, err := r.db.Query(`
		SELECT id, nim_sengketa, email_pelapor, path_foto_ktm, status_proses
		FROM sengketa_nim
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var disputes []domain.Dispute
	for rows.Next() {
		var d domain.Dispute
		if err := rows.Scan(&d.ID, &d.NIM, &d.ReporterEmail, &d.KTMPath, &d.Status); err != nil {
			return nil, err
		}
		disputes = append(disputes, d)
	}
	return disputes, nil
}

func (r *disputeRepository) ApproveAndResolveTransaction(disputeID int64) error {
	// Mulai Database Transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	// Pastikan rollback otomatis jika terjadi panic atau error sebelum di-commit
	defer tx.Rollback()

	var nim, reporterEmail string
	err = tx.QueryRow("SELECT nim, reporter_email FROM sengketa_nim WHERE id = $1", disputeID).Scan(&nim, &reporterEmail)
	if err != nil {
		return errors.New("sengketa tidak ditemukan")
	}

	// 2. Ambil Email Pembajak (Fraud) dari tabel pemilih
	var fraudEmail sql.NullString
	err = tx.QueryRow("SELECT email_gmail_login FROM pemilih WHERE nim = $1", nim).Scan(&fraudEmail)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// 3. Masukkan email pembajak ke Blacklist (FR-08)
	if fraudEmail.Valid && fraudEmail.String != reporterEmail {
		_, err = tx.Exec(`
    INSERT INTO email_blacklist (email, reason, created_at) 
    VALUES ($1, $2, NOW()) 
    ON CONFLICT (email) DO NOTHING
`, fraudEmail.String, "Terbukti membajak NIM "+nim)
		if err != nil {
			return err
		}
	}

	secret := os.Getenv("VOTE_SECRET_KEY")
	if secret == "" {
		// Fallback kalau di .env belum diset
		secret = "kunci-rahasia-pemira"
	}

	// 2. Panggil fungsi yang bener dari utils lu
	nimHash := utils.GenerateVoteHash(nim, secret)

	// 4. Hapus kertas suara lama milik pembajak (FR-08)
	_, err = tx.Exec("DELETE FROM kertas_suara WHERE voter_hash = $1", nimHash)
	if err != nil {
		return err
	}

	// 5. Reset data pemilih (Kosongkan email lama, hapus status milih & suspend) (FR-08)
	_, err = tx.Exec(`
		UPDATE pemilih 
		SET email_gmail_login = NULL, status_memilih = false, is_suspended = false 
		WHERE nim = $1
	`, nim)
	if err != nil {
		return err
	}

	// 6. Update status sengketa menjadi 'approved'
	_, err = tx.Exec("UPDATE sengketa_nim SET status_proses = 'approved' WHERE id = $1", disputeID)
	if err != nil {
		return err
	}

	// 7. Catat Audit Log (FR-08)
	_, err = tx.Exec(`
		INSERT INTO voter_events (voter_hash, event_type, description, created_at) 
		VALUES ($1, 'DISPUTE_APPROVED', 'Sengketa disetujui, suara lama dihapus, hak suara direset', NOW())
	`, nimHash)
	if err != nil {
		return err
	}

	return tx.Commit()
}
