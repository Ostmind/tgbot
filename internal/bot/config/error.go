package config

import "errors"

var (
	ErrNoBotToken   = errors.New("no tgbot token provided")
	ErrNoDBHost     = errors.New("no DB host provided")
	ErrNoDBPort     = errors.New("no DB port provided")
	ErrNoDBName     = errors.New("no DB name provided")
	ErrNoDBUser     = errors.New("no DB user provided")
	ErrNoDBPassword = errors.New("no DB password provided")
)
