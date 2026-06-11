package repository

import (
	"database/sql"
	"pemira-rpl/internal/domain"
)

type electionRepository struct{}

type ElectionRepository interface {
	GetByID(
		tx *sql.Tx,
		id int,
	) (*domain.Election, error)
}

func NewElectionRepository() ElectionRepository {
	return &electionRepository{}
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
