package service

import (
	"pemira-rpl/internal/domain"
	"pemira-rpl/internal/repository"
)

type CandidateService interface {
	GetAll() ([]domain.Candidate, error)
	GetResults() ([]domain.Result, error)
}

type candidateService struct {
	repo repository.CandidateRepository
}

func NewCandidateService(
	repo repository.CandidateRepository,
) CandidateService {
	return &candidateService{
		repo: repo,
	}
}

func (s *candidateService) GetAll() (
	[]domain.Candidate,
	error,
) {
	return s.repo.GetAll()
}

func (s *candidateService) GetResults() (
	[]domain.Result,
	error,
) {
	return s.repo.GetResults()
}
