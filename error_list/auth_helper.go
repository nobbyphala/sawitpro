package error_list

import "errors"

var (
	ErrPasswordNotMatch = errors.New("error password not match with hashed password")
	ErrInvalidToken     = errors.New("error invalid token")
	ErrNotAuthenticated = errors.New("error not authenticated")
)
