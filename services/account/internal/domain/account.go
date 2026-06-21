package domain

import (
	"time"
)

type Account struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	AccountNo   int64     `json:"account_no"`
	Balance     int64     `json:"balance"`
	AccountType ACCType   `json:"account_type"`
	Currency    string    `json:"currency"`
	Status      AccStatus `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
