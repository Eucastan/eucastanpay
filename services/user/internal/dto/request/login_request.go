package request

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
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
