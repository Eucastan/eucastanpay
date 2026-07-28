package account

type GetBalanceRequest struct {
	AccountNo int64 `json:"account_no" binding:"omitempty"`
}

type ActionRequest struct {
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
