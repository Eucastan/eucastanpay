package request

type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Phone       string `json:"phone" binding:"omitempty"`
	Password    string `json:"password_hash" binding:"required,min=8"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	DateOfBirth string `json:"date_of_birth" binding:"omitempty"`
}

type UpdateRequest struct {
	Email         string `json:"email" binding:"omitempty,email"`
	Phone         string `json:"phone" binding:"omitempty"`
	Password      string `json:"password_hash" binding:"omitempty,min=8"`
	FirstName     string `json:"first_name" binding:"omitempty"`
	LastName      string `json:"last_name" binding:"omitempty"`
	Status        string `json:"status" binding:"omitempty"`
	EmailVerified bool   `json:"email_verified" binding:"omitempty"`
}

type CreateKYCRequest struct {
	IdNumber string `json:"id_number" binding:"required,min=11"`
	IdType   string `json:"id_type" binding:"required,oneof=NIN passport license"`
}
