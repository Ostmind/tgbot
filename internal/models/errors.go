package models

import (
	"errors"
)

var (
	ErrUnique               = errors.New("already exists")
	ErrNotFound             = errors.New("not found")
	ErrDB                   = errors.New("db error")
	ErrDBConnectionCreation = errors.New("db connection creation error")
	ErrConfigCreation       = errors.New("config creation error")
	ErrNoServerHost         = errors.New("no server host provided")
	ErrNoServerPort         = errors.New("no server port provided")
	ErrNoDBHost             = errors.New("no DB host provided")
	ErrNoDBPort             = errors.New("no DB port provided")
	ErrNoDBName             = errors.New("no DB name provided")
	ErrNoDBUser             = errors.New("no DB user provided")
	ErrNoDBPassword         = errors.New("no DB password provided")
	ErrNoJWTSecret          = errors.New("no JWT Secret provided")
)
