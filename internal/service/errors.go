package service

import "errors"

// Domain errors.
var (
	ErrNotFound           = errors.New("resource not found")
	ErrForbidden          = errors.New("access denied")
	ErrAlreadyBooked      = errors.New("item already booked")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserExists         = errors.New("user already exists")
)
