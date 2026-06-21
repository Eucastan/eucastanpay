package domain

import (
	"time"
)

type KYCStatus string

const (
	StatusKycPending KYCStatus = "pending"
	StatusApproved   KYCStatus = "approved"
	StatusRejected   KYCStatus = "rejected"
)

type KYC struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	IDType     string     `json:"id_type"` // NIN, Passport, Drivers license
	IDNumber   string     `json:"id_number"`
	Status     KYCStatus  `json:"status"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}
