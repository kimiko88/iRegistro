package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidPIN         = errors.New("invalid security PIN")
	ErrNotFound           = errors.New("record not found")
)
