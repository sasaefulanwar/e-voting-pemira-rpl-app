package service

import (
	"database/sql"
	"errors"
	"fmt"
)

type VoteService interface {
	CastVote(voterID int, candidateID int) error
}

type voteService struct {
	db *sql.DB
}

func NewVoteService(db *sql.DB) VoteService {
	return &voteService{db: db}
}

func (s *voteService) CastVote(voterID int, candidateID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var hasVoted bool
	err = tx.QueryRow("SELECT has_voted FROM voters WHERE id = $1", voterID).Scan(&hasVoted)
	if err != nil {
		return fmt.Errorf("gagal cek status voter: %v", err)
	}

	if hasVoted {
		return errors.New("lu udah nyoblos cuy, ga bisa 2 kali!")
	}

	_, err = tx.Exec("INSERT INTO votes (voter_id, candidate_id) VALUES ($1, $2)", voterID, candidateID)
	if err != nil {
		return fmt.Errorf("gagal nyimpen suara: %v", err)
	}

	_, err = tx.Exec("UPDATE voters SET has_voted = TRUE WHERE id = $1", voterID)
	if err != nil {
		return fmt.Errorf("gagal update status voter: %v", err)
	}

	// 4. Sah!
	return tx.Commit()
}
