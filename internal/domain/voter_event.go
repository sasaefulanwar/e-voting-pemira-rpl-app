package domain

import "time"

type VoterEvent struct {
	ID        int64
	HashedNIM string
	EventType string
	IPAddress string
	UserAgent string
	CreatedAt time.Time
}
