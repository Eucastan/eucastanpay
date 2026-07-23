package request

type AccountRequest struct {
	Status    string `json:"status" binding:"required,oneof=active freeze closed"`
	Reason    string `json:"reason" binding:"required"`
	AccountNo int64  `json:"account_no" binding:"required"`
}

type UserRequest struct {
	Status    string `json:"status" binding:"required,oneof=active freeze closed"`
	Reason    string `json:"reason" binding:"required"`
	AccountNo int64  `json:"account_no" binding:"required"`
}
