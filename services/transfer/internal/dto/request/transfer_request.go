package request

type TransferRequest struct {
	FromAccID   string `json:"from_account_id" binding:"required"`
	FromAccNo   int64  `json:"from_account_no" binding:"required"`
	ToAccID     string `json:"to_account_id,omitempty"` // Inter-bank transfer
	ToAccNo     int64  `json:"to_account_no" binding:"required"`
	Amount      int64  `json:"amount" binding:"required,gt=0"`
	Description string `json:"description"`
	Mode        string `json:"mode" binding:"required,oneof=intraBank interBank own"`
}
