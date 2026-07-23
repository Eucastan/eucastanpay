package ledger

type LedgerRequest struct {
	Amount    int64  `json:"amount" binding:"required"`
	EntryType string `json:"entry_type" binding:"required"`
	Reference string `json:"reference" binding:"required"`
}

type EntryTypeRequest struct {
	EntryType string `json:"entry_type,omitempty"`
}

type EntryTypeURI struct {
	EntryType string `uri:"entry_type" binding:"omitempty"`
}

type AccountURI struct {
	AccountID string `uri:"account_id" binding:"required"`
}

type LedgerURI struct {
	LedgerID string `uri:"id" binding:"required"`
}

type Pagination struct {
	Limit int `form:"limit,default=10"`

	Page int `form:"page,default=1"`
}
