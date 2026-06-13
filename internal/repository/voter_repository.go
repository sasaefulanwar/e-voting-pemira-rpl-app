package repository

import (
	"database/sql"
	"errors"
	"pemira-rpl/internal/domain"
)

type VoterRepository interface {
	GetByNIMForUpdate(tx *sql.Tx, nim string) (*domain.Pemilih, error)
	UpdateEmail(tx *sql.Tx, nim string, email string) error
	UpdateStatusMemilih(
		tx *sql.Tx,
		nim string,
		status bool,
	) error
	FindByEmail(
		tx *sql.Tx,
		email string,
	) (*domain.Pemilih, error)
	GetStatistics(
		db *sql.DB,
	) (
		total int,
		voted int,
		notVoted int,
		err error,
	)
	SuspendByNIM(
		nim string,
	) error
	UnsuspendByNIM(
		nim string,
	) error
	IsEmailBlacklisted(email string) (bool, error)
}

type voterRepository struct {
	db *sql.DB
}

func NewVoterRepository(
	db *sql.DB,
) VoterRepository {

	return &voterRepository{
		db: db,
	}
}

func (r *voterRepository) GetByNIMForUpdate(tx *sql.Tx, nim string) (*domain.Pemilih, error) {
	var v domain.Pemilih
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

func (r *voterRepository) UpdateStatusMemilih(
	tx *sql.Tx,
	nim string,
	status bool,
) error {

	query := `
		UPDATE pemilih
		SET status_memilih = $1
		WHERE nim = $2
	`

	_, err := tx.Exec(
		query,
		status,
		nim,
	)

	return err
}

func (r *voterRepository) UpdateEmail(
	tx *sql.Tx,
	nim string,
	email string,
) error {

	query := `
		UPDATE pemilih
		SET email_gmail_login = $1
		WHERE nim = $2
	`

	_, err := tx.Exec(
		query,
		email,
		nim,
	)

	return err
}

func (r *voterRepository) FindByEmail(
	tx *sql.Tx,
	email string,
) (*domain.Pemilih, error) {

	var voter domain.Pemilih

	query := `
	SELECT
		nim,
		nama,
		email_gmail_login,
		status_memilih,
		is_suspended
	FROM pemilih
	WHERE email_gmail_login = $1
	LIMIT 1
	`

	err := tx.QueryRow(
		query,
		email,
	).Scan(
		&voter.NIM,
		&voter.Nama,
		&voter.EmailGmailLogin,
		&voter.StatusMemilih,
		&voter.IsSuspended,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &voter, nil
}

func (r *voterRepository) GetStatistics(
	db *sql.DB,
) (
	int,
	int,
	int,
	error,
) {

	var total int
	var voted int

	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM pemilih
	`).Scan(&total)

	if err != nil {
		return 0, 0, 0, err
	}

	err = db.QueryRow(`
		SELECT COUNT(*)
		FROM pemilih
		WHERE status_memilih = TRUE
	`).Scan(&voted)

	if err != nil {
		return 0, 0, 0, err
	}

	return total, voted, total - voted, nil
}

func (r *voterRepository) SuspendByNIM(nim string) error {
	res, err := r.db.Exec(`
        UPDATE pemilih
        SET is_suspended = TRUE
        WHERE nim = $1
    `, nim)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("NIM tidak ditemukan")
	}

	return nil
}

func (r *voterRepository) UnsuspendByNIM(nim string) error {
	_, err := r.db.Exec(`
		UPDATE pemilih
		SET is_suspended=FALSE
		WHERE nim=$1
	`, nim)
	return err
}

func (r *voterRepository) IsEmailBlacklisted(email string) (bool, error) {
	var isBlacklisted bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM email_blacklist WHERE email = $1)", email).Scan(&isBlacklisted)
	if err != nil {
		return false, err
	}
	return isBlacklisted, nil
}
