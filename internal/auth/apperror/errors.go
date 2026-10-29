package apperror

import "errors"

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
