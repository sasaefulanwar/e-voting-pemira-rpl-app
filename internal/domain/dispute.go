package domain

import "time"

type Dispute struct {
	ID int64 `json:"id"`

	NIM string `json:"nim"`

	ReporterEmail string `json:"reporter_email"`

	KTMPath string `json:"ktm_path"`

	Status string `json:"status"`

	CreatedAt time.Time `json:"created_at"`
}
