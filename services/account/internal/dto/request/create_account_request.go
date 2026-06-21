package request

type CreateAccountRequest struct {
	AccountNo   int64  `json:"account_no" binding:"required"`
	Balance     int64  `json:"balance" binding:"required"`
	AccountType string `json:"account_type" binding:"required,oneof=current fixed_deposit savings"`
	Currency    string `json:"currency" binding:"required,oneof=NGN USD EURO"`
}
