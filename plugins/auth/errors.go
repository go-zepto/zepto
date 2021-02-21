package auth

import "errors"

var ErrInvalidToken = errors.New("invalid token")
var ErrInvalidPID = errors.New("invalid PID")
var ErrUnauthorized = errors.New("unauthorized")
