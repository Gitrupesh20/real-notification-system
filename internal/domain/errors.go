// internal/domain/errors.go
package domain

import "errors"

// Domain-specific errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidUserID     = errors.New("invalid user ID")
	ErrConnectionClosed  = errors.New("connection closed")
	ErrMessageInvalid    = errors.New("message validation failed")
	ErrUserAlreadyExists = errors.New("user already exists")

	//Message
	ErrMessageIsNil = errors.New("message is nil")
)
