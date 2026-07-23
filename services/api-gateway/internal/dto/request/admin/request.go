package admin

type CreateAdminRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role" binding:"required"`
}

type UpdateAdminRequest struct {
	Email     *string `json:"email" binding:"omitempty"`
	Password  *string `json:"password" binding:"omitempty,min=8"`
	FirstName *string `json:"first_name" binding:"omitempty"`
	LastName  *string `json:"last_name" binding:"omitempty"`
	Role      *string `json:"role" binding:"omitempty"`
	Status    *string `json:"status" binding:"omitempty"`
}

type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	TotpCode string `json:"totp_code" binding:"required"`
}

type ReverseTransactionRequest struct {
	Reference string `json:"reference" binding:"required"`
	Reason    string `json:"reason"`
}

type ReconciliationRequest struct {
	AccountNo int64 `json:"account_no" binding:"required"`
}

type AdminURI struct {
	AdminID string `uri:"id" binding:"required"`
}

type Pagination struct {
	Limit int `form:"limit,default=10"`
	Page  int `form:"page,default=1"`
}
