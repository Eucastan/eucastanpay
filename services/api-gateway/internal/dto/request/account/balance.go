package account

type GetBalanceRequest struct {
	AccountID string `json:"account_id" binding:"required"`
	AccountNo int64  `json:"account_no" binding:"omitempty"`
	UserID    string `json:"user_id" binding:"required"`
}

type ActionRequest struct {
	AccountID string `json:"account_id" binding:"required"`
	Status    string `json:"status" binding:"omitempty"`
	AccountNo int64  `json:"account_no" binding:"omitempty"`
}

type AccountURI struct {
	AccountID string `uri:"id" binding:"required"`
}

type Pagination struct {
	Limit int `form:"limit,default=10"`

	Page int `form:"page,default=1"`
}
