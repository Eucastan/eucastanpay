package response

type AuthResponse struct {
	Message      string       `json:"message"`
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

type RegisterResponse struct {
	Message string       `json:"message"`
	User    UserResponse `json:"user"`
}

type LoginResponse struct {
	Message  string       `json:"message"`
	Response AuthResponse `json:"response"`
}

type KYCResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
