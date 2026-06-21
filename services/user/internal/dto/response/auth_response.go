package response

type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

type KYCResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
