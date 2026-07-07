package errmessage

import (
	"errors"
)

var (
	ErrTranferNotFound                    = errors.New("transfer not found")
	ErrNothingUpdated                     = errors.New("no rows updated")
	ErrPreviousTransferFailure            = errors.New("previous transfer failed, please retry with new key")
	ErrCanNotTransferToSelf               = errors.New("cannot transfer to self")
	ErrCannotReverseNonSuccessfulTransfer = errors.New("cannot reverse non successful transfer")
	ErrAlreadyReversed                    = errors.New("already reversed")
	ErrUserNotOwner                       = errors.New("user does not own account")
	ErrInvalidTransferMode                = errors.New("invalid transfer mode")
	ErrDuplicateRequest                   = errors.New("duplicate request")
	ErrUnauthorized                       = errors.New("unauthorized")
	ErrLedgerNotFound                     = errors.New("ledger not found")
)
