package account

import "time"

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type AccountResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Email       string    `json:"email"`
	AccountNo   int64     `json:"account_no"`
	Balance     int64     `json:"amount"`
	AccountType string    `json:"account_type"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ConfirmAccountResponse struct {
	FromAccID   string `json:"from_account_id"`
	ToAccID     string `json:"to_account_id"`
	FromUserID  string `json:"sender"`
	ToUserID    string `json:"receiver"`
	FromEmail   string `json:"from_email"`
	ToEmail     string `json:"to_email"`
	FromBalance int64  `json:"from_balance"`
	ToBalance   int64  `json:"to_balance"`
	FromStatus  string `json:"from_status"`
	ToStatus    string `json:"to_status"`
}
