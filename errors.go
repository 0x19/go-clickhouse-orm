package chorm

import "errors"

var (
	ErrNoConfigProvided   = errors.New("configuration must be provided")
	ErrNoHostProvided     = errors.New("host must be provided")
	ErrNoPortProvided     = errors.New("port must be provided")
	ErrNoUsernameProvided = errors.New("username must be provided")
	ErrNoPasswordProvided = errors.New("password must be provided")
	ErrNoDatabaseProvided = errors.New("database must be provided")
)
