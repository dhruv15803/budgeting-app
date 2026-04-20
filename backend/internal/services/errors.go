package services

import "errors"

var (
	ErrEmailTaken             = errors.New("email already registered")
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrNotVerified            = errors.New("email address is not verified")
	ErrInvalidOrExpiredToken  = errors.New("invalid or expired verification token")
	ErrNotFound               = errors.New("resource not found")
	ErrForbidden              = errors.New("you do not have permission to perform this action")
	ErrGoogleAuthDisabled     = errors.New("Google sign-in is not configured")
	ErrInvalidGoogleToken     = errors.New("invalid Google credential")
	ErrGoogleEmailNotVerified = errors.New("Google account email is not verified")
	ErrGoogleAccountConflict  = errors.New("this email is linked to a different Google account")
)
