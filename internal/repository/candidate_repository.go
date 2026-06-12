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
	rows, err := r.db.Query(`
	SELECT
		id_paslon,
		election_id,

		nama_ketua,
		nama_wakil,

		visi,
		misi,

		foto_paslon

	FROM kandidat
`)
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

			&c.ChairmanName,
			&c.ViceChairmanName,

			&c.Vision,
			&c.Mission,

			&c.PhotoURL,
		); err != nil {
			return nil, err
		}
		candidates = append(candidates, c)
	}

	return candidates, nil
}

func (r *candidateRepository) GetResults() ([]domain.Result, error) {
	// Slice literal supaya JSON encode selalu []
	results := []domain.Result{}

	// LEFT JOIN ke kandidat supaya paslon tanpa suara tetap muncul
	rows, err := r.db.Query(`
		SELECT
			k.id AS id_paslon,
			COALESCE(COUNT(s.id_paslon), 0) AS votes
		FROM kandidat k
		LEFT JOIN kertas_suara s
			ON k.id = s.id_paslon
		GROUP BY k.id
		ORDER BY k.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var result domain.Result
		if err := rows.Scan(&result.PaslonID, &result.Votes); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}
