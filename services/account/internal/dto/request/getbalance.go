package request

type GetBalanceRequest struct {
	AccountNo int64 `json:"account_no" binding:"required"`
}
