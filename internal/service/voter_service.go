package service

import (
	"database/sql"
	"errors"
	"pemira-rpl/internal/dto"
	"pemira-rpl/internal/repository"
)

type VoterService interface {
	ProcessBinding(req dto.BindNIMRequest, loggedInEmail string) (*dto.BindNIMResponse, error)
}

type voterService struct {
	db   *sql.DB
	repo repository.VoterRepository
}

func NewVoterService(db *sql.DB, repo repository.VoterRepository) VoterService {
	return &voterService{db: db, repo: repo}
}

func (s *voterService) ProcessBinding(req dto.BindNIMRequest, loggedInEmail string) (*dto.BindNIMResponse, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, errors.New("gagal memulai transaksi server")
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	voter, err := s.repo.GetByNIMForUpdate(tx, req.NIM)
	if err != nil {
		return nil, err
	}

	if voter.IsSuspended {
		err = errors.New("NIM kamu ditangguhkan (Suspended)")
		return nil, err
	}
	if voter.EmailGmailLogin != "" {
		err = errors.New("NIM ini sudah di-binding dengan email lain")
		return nil, err
	}

	err = s.repo.UpdateEmail(tx, req.NIM, loggedInEmail)
	if err != nil {
		err = errors.New("gagal menyimpan email ke database")
		return nil, err
	}

	return &dto.BindNIMResponse{
		Message: "Mantap! NIM berhasil di-binding ke email Google kamu.",
		Status:  "success",
	}, nil
}
