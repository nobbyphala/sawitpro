package error_list

import "errors"

var (
	ErrProfileRegister = errors.New("error when register a new profile")
	ErrGetProfile      = errors.New("error when get user profile")
	ErrProfileNotFound = errors.New("error profile not found")
	ErrDataConflict    = errors.New("error there existing data conficted with new data")

	ErrLoginCredential = errors.New("error credentials combination not match")
	ErrLogin           = errors.New("error when try to login")

	ErrUpdateProfile = errors.New("error when updating profile")
)
