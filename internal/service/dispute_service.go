package service

import (
	"pemira-rpl/internal/domain"
	"pemira-rpl/internal/repository"
)

type DisputeService interface {
	SubmitDispute(nim, reporterEmail, ktmPath string) error
	ApproveDispute(id int64) error
	RejectDispute(id int64) error
	GetAllDisputes() ([]domain.Dispute, error)
}

type disputeService struct {
	repo      repository.DisputeRepository
	voterRepo repository.VoterRepository
}

func NewDisputeService(repo repository.DisputeRepository, voterRepo repository.VoterRepository) DisputeService {
	return &disputeService{repo: repo, voterRepo: voterRepo}
}

func (s *disputeService) SubmitDispute(nim, reporterEmail, ktmPath string) error {
	err := s.repo.Submit(domain.Dispute{
		NIM:           nim,
		ReporterEmail: reporterEmail,
		KTMPath:       ktmPath,
		Status:        "pending",
	})
	if err != nil {
		return err
	}

	err = s.voterRepo.SuspendByNIM(nim)
	if err != nil {
		return err
	}

	return nil
}

func (s *disputeService) ApproveDispute(disputeID int64) error { // <-- Ganti jadi int64
	err := s.repo.ApproveAndResolveTransaction(disputeID)
	if err != nil {
		return err
	}
	return nil
}

func (s *disputeService) RejectDispute(id int64) error {
	return s.repo.Reject(id)
}

func (s *disputeService) GetAllDisputes() ([]domain.Dispute, error) {
	return s.repo.GetAll()
}
