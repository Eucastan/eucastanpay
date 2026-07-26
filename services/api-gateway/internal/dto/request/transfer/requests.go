package transfer

type TransferRequest struct {
	UserID      string `json:"user_id" binding:"required"`
	IdemKey     string `json:"idempotency_key" binding:"required"`
	ToAccNo     int64  `json:"to_account_no" binding:"required"`
	Amount      int64  `json:"amount" binding:"required,gt=0"`
	Description string `json:"description"`
	Mode        string `json:"mode" binding:"required,oneof=intraBank interBank own"`
}

type ReverseTransferRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	Reference string `json:"reference" binding:"required"`
	IdemKey   string `json:"idempotency_key" binding:"required"`
	ToAccNo   int64  `json:"to_account_no" binding:"required"`
	Amount    int64  `json:"amount" binding:"required,gt=0"`
}

type ReconciliationRequest struct {
	AccountNo int64 `json:"to_account_no" binding:"required"`
}

type TransactionIdentity struct {
	FromAccID   string `json:"from_account_id" binding:"required"`
	FromAccNo   int64  `json:"from_account_no" binding:"required"`
	ToAccID     string `json:"to_account_id" binding:"required"`
	ToAccNo     int64  `json:"to_account_no" binding:"required"`
	Amount      int64  `json:"amount" binding:"required,gt=0"`
	Description string `json:"description"`
	Mode        string `json:"mode" binding:"required,oneof=intraBank interBank own"`
}

type TransferURI struct {
	TranferID string `uri:"id" binding:"required"`
}

type ReverseURI struct {
	Reference string `uri:"reference" binding:"required"`
}

type AccountIdURI struct {
	AccountID string `uri:"account_id" binding:"required"`
}

type Pagination struct {
	Limit int `form:"limit,default=10"`

	Page int `form:"page,default=1"`
}
