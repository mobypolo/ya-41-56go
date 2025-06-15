package errors

import "errors"

var (
	ErrInitDB         = errors.New("failed to init database")
	ErrEmptySecretKey = errors.New("empty secret key for JWT")
	ErrHTTPServer     = errors.New("failed to start HTTP server")

	ErrInvalidCreds = errors.New("invalid credentials")
	ErrUserExists   = errors.New("user already exists")
	ErrJWTToken     = errors.New("invalid JWT token")
)
