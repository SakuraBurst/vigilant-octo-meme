package services

import "github.com/go-faster/errors"

var (
	ErrUserNotAuthorized  = errors.New("user not authorize")
	ErrUserDontHaveAccess = errors.New("user don't have access")
	ErrTokenInvalid       = errors.New("token invalid")
)
