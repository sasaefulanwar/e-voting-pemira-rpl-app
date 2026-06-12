package service

import (
	"errors"
	"pemira-rpl/internal/domain"
	"pemira-rpl/internal/repository"
)

type CandidateService interface {
	GetAll() ([]domain.Candidate, error)

	GetResults(
		isAdmin bool,
	) ([]domain.Result, error)
}

type candidateService struct {
	repo         repository.CandidateRepository
	electionRepo repository.ElectionRepository
}

func NewCandidateService(
	repo repository.CandidateRepository,
	electionRepo repository.ElectionRepository,
) CandidateService {

	return &candidateService{
		repo:         repo,
		electionRepo: electionRepo,
	}
}

func (s *candidateService) GetAll() (
	[]domain.Candidate,
	error,
) {
	return s.repo.GetAll()
}

func (s *candidateService) GetResults(
	isAdmin bool,
) ([]domain.Result, error) {

	election, err :=
		s.electionRepo.GetCurrent()

	if err != nil {
		return nil, err
	}

	if election.Status == "open" &&
		!isAdmin {

		return nil,
			errors.New(
				"hasil belum dapat diakses",
			)
	}

	return s.repo.GetResults()
}
