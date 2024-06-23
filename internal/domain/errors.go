package domain

import "errors"

var (
	ErrUnauthorized   = errors.New("unauthorized access")
	ErrInternalServer = errors.New("internal server error")
	ErrBadRequest     = errors.New("bad request")
)
