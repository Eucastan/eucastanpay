package domain

import (
	"time"
)

type Account struct {
	ID          string    `db:"id" json:"id"`
	UserID      string    `db:"user_id" json:"user_id"`
	Email       string    `db:"email" json:"email"`
	AccountNo   int64     `db:"account_no" json:"account_no"`
	Balance     int64     `db:"balance" json:"balance"`
	AccountType ACCType   `db:"account_type" json:"account_type"`
	Currency    string    `db:"currency" json:"currency"`
	Status      AccStatus `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
