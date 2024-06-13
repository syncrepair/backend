package domain

import "errors"

var (
	ErrInternalServer = errors.New("внутрішня помилка сервера")
	ErrBadRequest     = errors.New("невірний запит")
	ErrUnauthorized   = errors.New("помилка авторизації")
)
