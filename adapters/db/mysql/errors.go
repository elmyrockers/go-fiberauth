package mysql

import "errors"

var (
	ErrNotFound       = errors.New("mysql: record not found")
	ErrDuplicateEmail = errors.New("mysql: email already registered")
)