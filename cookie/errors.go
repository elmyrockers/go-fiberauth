package cookie

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailTaken       = errors.New("fiberauth: email already registered")
	ErrNotFound         = errors.New("fiberauth: record not found")
	ErrInvalidCredentials = errors.New("fiberauth: invalid credentials")
	ErrTokenExpired     = errors.New("fiberauth: token expired")
	ErrTokenUsed        = errors.New("fiberauth: token already used")
	ErrInvalidToken 	= errors.New("fiberauth: invalid or expired token")

	ErrPasswordTooLong = bcrypt.ErrPasswordTooLong
)