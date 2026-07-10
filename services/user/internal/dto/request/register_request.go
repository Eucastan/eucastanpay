package request

type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Phone       string `json:"phone" binding:"omitempty"`
	Password    string `json:"password_hash" binding:"required,min=8"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	DateOfBirth string `json:"date_of_birth" binding:"omitempty"`
}
