package services

import "errors"

var (
	ErrEmailTaken            = errors.New("email already registered")
	ErrInvalidCredentials    = errors.New("invalid email or password")
	ErrNotVerified           = errors.New("email address is not verified")
	ErrInvalidOrExpiredToken = errors.New("invalid or expired verification token")
)
