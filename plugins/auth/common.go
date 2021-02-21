package auth

import "time"

type PID interface{}

type Token struct {
	Value      string     `json:"value"`
	Expiration *time.Time `json:"expiration"`
}

type AuthCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
