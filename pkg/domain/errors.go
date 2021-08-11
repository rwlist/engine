package domain

import "errors"

var (
	ErrDatabaseNotFound   = errors.New("database not found")
	ErrDatabaseLoadFailed = errors.New("failed to load database")
	ErrDatabaseExists     = errors.New("database already exists")
	ErrAccessDenied       = errors.New("access denied")
	ErrInvalidList        = errors.New("list is invalid")
	ErrListNotFound       = errors.New("list not found")
)
