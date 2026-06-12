package service

import (
	"database/sql"
	"errors"
	"os"
	"pemira-rpl/internal/domain"
	"pemira-rpl/internal/repository"
	"pemira-rpl/internal/utils"
)

type VoteService interface {
	CastVote(
		nim string,
		electionID int,
		paslonID int,
	) error
}

type voteService struct {
	db           *sql.DB
	voterRepo    repository.VoterRepository
	voteRepo     repository.VoteRepository
	electionRepo repository.ElectionRepository
	auditRepo    repository.AuditRepository
}

func NewVoteService(
	db *sql.DB,
	voterRepo repository.VoterRepository,
	voteRepo repository.VoteRepository,
	electionRepo repository.ElectionRepository,
	auditRepo repository.AuditRepository,
) VoteService {

	return &voteService{
		db:           db,
		voterRepo:    voterRepo,
		voteRepo:     voteRepo,
		electionRepo: electionRepo,
		auditRepo:    auditRepo,
	}
}

func (s *voteService) CastVote(
	nim string,
	electionID int,
	paslonID int,
) error {

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	voter, err := s.voterRepo.
		GetByNIMForUpdate(tx, nim)

	if err != nil {
		return err
	}

	if voter.IsSuspended {
		return errors.New(
			"akun sedang disuspend",
		)
	}

	if voter.StatusMemilih {
		return errors.New(
			"sudah memilih",
		)
	}

	election, err := s.electionRepo.
		GetByID(tx, electionID)

	if err != nil {
		return err
	}

	if election.Status != "open" {
		return errors.New(
			"pemilu belum dibuka",
		)
	}

	hashedNIM :=
		utils.GenerateVoteHash(
			nim,
			os.Getenv("VOTE_SECRET_KEY"),
		)

	ballot := &domain.Ballot{
		ElectionID: electionID,
		HashedNIM:  hashedNIM,
		PaslonID:   paslonID,
	}

	err = s.voteRepo.
		InsertBallot(tx, ballot)

	if err != nil {
		return err
	}

	err = s.voterRepo.
		UpdateStatusMemilih(
			tx,
			nim,
			true,
		)

	if err != nil {
		return err
	}

	err = s.auditRepo.LogEventTx(
		tx,
		hashedNIM,
		"VOTE_CAST",
	)

	if err != nil {
		return err
	}

	return tx.Commit()
}
