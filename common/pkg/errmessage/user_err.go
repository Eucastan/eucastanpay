package errmessage

import (
	"errors"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrDuplicateEmail       = errors.New("email already exists")
	ErrKYCAlreadyExists     = errors.New("kyc already exists")
	ErrKYCAlreadySubmitted  = errors.New("kyc already submitted")
	ErrInvalidToken         = errors.New("invalid token")
	ErrExpiredToken         = errors.New("expired token")
	ErrPasswordNotConfirmed = errors.New("password mismatch, try again")
)
