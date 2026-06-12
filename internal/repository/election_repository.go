package repository

import (
	"database/sql"
	"pemira-rpl/internal/domain"
)

type electionRepository struct {
	db *sql.DB
}

type ElectionRepository interface {
	GetByID(
		tx *sql.Tx,
		id int,
	) (*domain.Election, error)

	UpdateStatus(
		status string,
	) error

	GetCurrent() (*domain.Election, error)
}

func NewElectionRepository(db *sql.DB) ElectionRepository {
	return &electionRepository{db: db}
}

func (r *electionRepository) GetByID(
	tx *sql.Tx,
	id int,
) (*domain.Election, error) {
	query := `
		SELECT id, nama_pemilu, start_at, end_at, status
		FROM elections
		WHERE id = $1
	`
	election := &domain.Election{}
	row := tx.QueryRow(query, id)
	if err := row.Scan(
		&election.ID,
		&election.NamaPemilu,
		&election.StartAt,
		&election.EndAt,
		&election.Status,
	); err != nil {
		return nil, err
	}
	return election, nil
}

func (r *electionRepository) UpdateStatus(
	status string,
) error {

	_, err := r.db.Exec(
		`
        UPDATE elections
        SET status = $1
        WHERE id = 1
        `,
		status,
	)

	return err
}

func (r *electionRepository) GetCurrent() (
	*domain.Election,
	error,
) {

	query := `
		SELECT
			id,
			nama_pemilu,
			start_at,
			end_at,
			status
		FROM elections
		WHERE id = 1
	`

	election := &domain.Election{}

	err := r.db.QueryRow(query).Scan(
		&election.ID,
		&election.NamaPemilu,
		&election.StartAt,
		&election.EndAt,
		&election.Status,
	)

	if err != nil {
		return nil, err
	}

	return election, nil
}
