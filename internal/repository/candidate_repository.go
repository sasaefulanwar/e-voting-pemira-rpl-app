package repository

import (
	"database/sql"
	"pemira-rpl/internal/domain"
)

type CandidateRepository interface {
	GetAll() ([]domain.Candidate, error)
	GetResults() ([]domain.Result, error)
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

func (r *candidateRepository) GetResults() (
	[]domain.Result,
	error,
) {

	rows, err := r.db.Query(`
		SELECT
			id_paslon,
			COUNT(*)
		FROM kertas_suara
		GROUP BY id_paslon
		ORDER BY id_paslon
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.Result

	for rows.Next() {

		var result domain.Result

		if err := rows.Scan(
			&result.PaslonID,
			&result.Votes,
		); err != nil {
			return nil, err
		}

		results = append(
			results,
			result,
		)
	}

	return results, nil
}
