package request

type CreateAccountRequest struct {
	AccountType string `json:"account_type" binding:"required,oneof=current fixed_deposit savings"`
	Currency    string `json:"currency" binding:"required,oneof=NGN USD EURO"`
}
