package errmessage

import (
	"errors"
)

var (
	ErrAccNotFound             = errors.New("not found")
	ErrAccountDisabled         = errors.New("account is disabled")
	ErrNotEligibleForOperation = errors.New("not eligible for this operation")
	ErrInsufficientAmount      = errors.New("insufficient amount")
	ErrDuplicateAccNo          = errors.New("duplicate account number")
	ErrAccAlreadyExists        = errors.New("account already exists")
	ErrIdemKeyNotFound         = errors.New("idempotencykey not found")
)
