package request

type LedgerRequest struct {
	Amount    int64  `json:"amount" binding:"required"`
	EntryType string `json:"entry_type" binding:"required"`
	Reference string `json:"reference" binding:"required"`
}

type EntryTypeRequest struct {
	EntryType string `json:"entry_type,omitempty"`
}
