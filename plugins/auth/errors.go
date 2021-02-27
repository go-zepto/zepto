package auth

import "errors"

var ErrInvalidToken = errors.New("invalid token")
var ErrInvalidPID = errors.New("invalid PID")
var ErrUnauthorized = errors.New("unauthorized")
var ErrInternalServerError = errors.New("internal server error")
var ErrMissingUsernameOrPassword = errors.New("missing username or password")
