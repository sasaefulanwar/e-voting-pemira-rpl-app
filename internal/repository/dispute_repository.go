package repository

import (
	"database/sql"
	"pemira-rpl/internal/domain"
)

type DisputeRepository interface {
	Submit(dispute domain.Dispute) error
	GetByID(id int64) (*domain.Dispute, error)
	Approve(id int64) error
	Reject(id int64) error
	GetAll() ([]domain.Dispute, error)
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
