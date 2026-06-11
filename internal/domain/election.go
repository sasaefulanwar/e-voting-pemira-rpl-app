package domain

import "time"

type Election struct {
	ID         int
	NamaPemilu string
	StartAt    time.Time
	EndAt      time.Time
	Status     string
}
