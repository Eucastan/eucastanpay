package errmessage

import "errors"

var (
	ErrForbidden       = errors.New("forbidden")
	ErrInvalidRole     = errors.New("invalid role")
	ErrInvalidStatus   = errors.New("invalid status update")
	ErrAdminNotFound   = errors.New("admin not found")
	Err2FANotEnabled   = errors.New("2factor authentication not enabled")
	ErrInvalid2FACode  = errors.New("invalid 2factor authentication")
	ErrBootstrapClosed = errors.New("bootstrap closed")
)
