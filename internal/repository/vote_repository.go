package repository

import (
	"database/sql"
	"pemira-rpl/internal/domain"
)

type voteRepository struct{}

type VoteRepository interface {
	InsertBallot(
		tx *sql.Tx,
		ballot *domain.Ballot,
	) error
}

func NewVoteRepository() VoteRepository {
	return &voteRepository{}
}

func (r *voteRepository) InsertBallot(
	tx *sql.Tx,
	ballot *domain.Ballot,
) error {

	query := `
		INSERT INTO kertas_suara
		(
			election_id,
			hashed_nim,
			id_paslon
		)
		VALUES ($1,$2,$3)
	`

	_, err := tx.Exec(
		query,
		ballot.ElectionID,
		ballot.HashedNIM,
		ballot.PaslonID,
	)

	return err
}
