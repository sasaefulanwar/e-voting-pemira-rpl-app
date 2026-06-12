package service

import (
	"database/sql"
	"pemira-rpl/internal/dto"
	"pemira-rpl/internal/repository"
)

type ElectionService interface {
	OpenElection() error
	CloseElection() error

	GetStatistics() (
		*dto.AdminStatisticsResponse,
		error,
	)
}

type electionService struct {
	db        *sql.DB
	repo      repository.ElectionRepository
	voterRepo repository.VoterRepository
}

func NewElectionService(
	db *sql.DB,
	repo repository.ElectionRepository,
	voterRepo repository.VoterRepository,
) ElectionService {

	return &electionService{
		db:        db,
		repo:      repo,
		voterRepo: voterRepo,
	}
}

func (s *electionService) OpenElection() error {
	return s.repo.UpdateStatus("open")
}

func (s *electionService) CloseElection() error {
	return s.repo.UpdateStatus("closed")
}

func (s *electionService) GetStatistics() (
	*dto.AdminStatisticsResponse,
	error,
) {

	total,
		voted,
		notVoted,
		err := s.voterRepo.GetStatistics(
		s.db,
	)

	if err != nil {
		return nil, err
	}

	election,
		err := s.repo.GetCurrent()

	if err != nil {
		return nil, err
	}

	return &dto.AdminStatisticsResponse{
		TotalVoters:    total,
		Voted:          voted,
		NotVoted:       notVoted,
		ElectionStatus: election.Status,
	}, nil
}
