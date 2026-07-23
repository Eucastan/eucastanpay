package request

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserURI struct {
	UserID string `uri:"id" binding:"required"`
}

type UserKYCURI struct {
	UserID string `uri:"user_id" binding:"required"`
}

type Pagination struct {
	Limit int `form:"limit,default=10"`

	Page int `form:"page,default=1"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	Password    string `json:"password" binding:"required,min=8"`
	ConfirmPass string `json:"confirm_pass" binding:"required,min=8"`
}

type CurrentStatusRequest struct {
	Status string `json:"status"`
}
