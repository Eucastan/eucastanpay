package request

type KYCRequest struct {
	IDType   string `json:"id_type" binding:"required,oneof=NIN passport drivers_license"`
	IDNumber string `json:"id_number" binding:"required"`
}
