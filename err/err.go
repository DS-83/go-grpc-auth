package err

import (
	"errors"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCred        = errors.New("invalid credentials")
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrDupKey             = errors.New("username already in use")
)
