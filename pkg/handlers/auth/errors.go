package auth

import "errors"

var (
	ErrInvalidJwt = errors.New("invalid jwt")
)
