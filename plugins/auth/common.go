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

type AuthDatasource interface {
	Auth(username string, password string) (PID, error)
	ResetPassword(pid PID, password string) error
}

type AuthEncoder interface {
	GenerateTokenFromPID(pid PID) (*Token, error)
}

type AuthStore interface {
	StoreAuthToken(token *Token, pid PID) error
	GetAuthTokenPID(token string) (PID, error)
	StoreResetPasswordToken(token *Token, pid PID) error
	GetResetPasswordTokenPID(token string) (PID, error)
}
