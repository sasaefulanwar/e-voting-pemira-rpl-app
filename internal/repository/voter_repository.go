package repository

import (
	"database/sql"
	"errors"
	"pemira-rpl/internal/domain"
)

type VoterRepository interface {
	// Pakai *sql.Tx biar masuk dalam satu alur transaksi
	GetByNIMForUpdate(tx *sql.Tx, nim string) (*domain.Pemilih, error)
	UpdateEmail(tx *sql.Tx, nim string, email string) error
}

type voterRepository struct{}

func NewVoterRepository() VoterRepository {
	return &voterRepository{}
}

func (r *voterRepository) GetByNIMForUpdate(tx *sql.Tx, nim string) (*domain.Pemilih, error) {
	var v domain.Pemilih
	// FOR UPDATE ini kunci rahasianya cuy! Baris ini bakal dikunci sampai transaksi beres.
	query := `SELECT nim, nama, email_gmail_login, status_memilih, is_suspended 
	          FROM pemilih WHERE nim = $1 FOR UPDATE`

	err := tx.QueryRow(query, nim).Scan(&v.NIM, &v.Nama, &v.EmailGmailLogin, &v.StatusMemilih, &v.IsSuspended)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("NIM tidak terdaftar di DPT")
		}
		return nil, err
	}
	return &v, nil
}

func (r *voterRepository) UpdateEmail(tx *sql.Tx, nim string, email string) error {
	query := `UPDATE pemilih SET email_gmail_login = $1 WHERE nim = $2`
	_, err := tx.Exec(query, email, nim)
	return err
}
