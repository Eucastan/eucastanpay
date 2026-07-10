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
	ID         string     `db:"id" json:"id"`
	UserID     string     `db:"user_id" json:"user_id"`
	IDType     string     `db:"id_type" json:"id_type"` // NIN, Passport, Drivers license
	IDNumber   string     `db:"id_number" json:"id_number"`
	Status     KYCStatus  `db:"status" json:"status"`
	VerifiedAt *time.Time `db:"verified_at" json:"verified_at,omitempty"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
}
