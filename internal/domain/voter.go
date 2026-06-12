package domain

import "time"

type Pemilih struct {
	NIM             string    `json:"nim"`
	Nama            string    `json:"nama"`
	EmailGmailLogin *string   `json:"email_gmail_login,omitempty"`
	StatusMemilih   bool      `json:"status_memilih"`
	IsSuspended     bool      `json:"is_suspended"`
	CreatedAt       time.Time `json:"created_at"`
}
