package account

type DepositRequest struct {
	AccountID string `json:"account_id" binding:"required"`
	AccountNo int64  `json:"account_no" binding:"required"`
	Amount    int64  `json:"amount" binding:"required,gt=0"`
	Currency  string `json:"currency"`
}
