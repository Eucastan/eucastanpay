package request

type CreateAdminRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
}

type UpdateAdminRequest struct {
	Email     *string `json:"email"`
	Password  *string `json:"password" binding:"omitempty,min=8"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Role      *string `json:"role"`
	Status    *string `json:"status"`
}

type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ReverseTransactionRequest struct {
	Reference string `json:"reference" binding:"required"`
	Reason    string `json:"reason"`
}

type ReconciliationRequest struct {
	AccountNo int64 `json:"account_no" binding:"required"`
}
