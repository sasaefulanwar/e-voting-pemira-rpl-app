package service

import (
	"pemira-rpl/internal/domain"
	"pemira-rpl/internal/repository"
)

type CandidateService interface {
	GetAll() ([]domain.Candidate, error)
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
