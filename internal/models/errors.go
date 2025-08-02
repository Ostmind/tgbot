package models

import (
	"errors"
)

var (
	ErrUnique               = errors.New("already exists")
	ErrNotFound             = errors.New("not found")
	ErrDBConnectionCreation = errors.New("db connection creation error")
	ErrPasswordNotSecured   = errors.New("password doesn't met requirements")
)
