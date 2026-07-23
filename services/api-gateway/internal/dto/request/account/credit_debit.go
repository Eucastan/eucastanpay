package account

type CreditRequest struct {
	AccountNo int64 `json:"account_no" binding:"required"`
	Amount    int64 `json:"amount" binding:"required,gt=0"`
}

type DebitRequest struct {
	AccountNo int64 `json:"account_no" binding:"required"`
	Amount    int64 `json:"amount" binding:"required,gt=0"`
}
