package repository

import (
	"database/sql"
	"pemira-rpl/internal/domain"
)

type CandidateRepository interface {
	GetAll() ([]domain.Candidate, error)
}

type candidateRepository struct {
	db *sql.DB
}

func NewCandidateRepository(db *sql.DB) CandidateRepository {
	return &candidateRepository{
		db: db,
	}
}

func (r *candidateRepository) GetAll() ([]domain.Candidate, error) {
	rows, err := r.db.Query(`SELECT id_paslon, election_id, nama_ketua, nama_wakil FROM kandidat`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var candidates []domain.Candidate
	for rows.Next() {
		var c domain.Candidate
		if err := rows.Scan(
			&c.ID,
			&c.ElectionID,
			&c.Name,
			&c.Mission,
		); err != nil {
			return nil, err
		}
		candidates = append(candidates, c)
	}

	return candidates, nil
}
