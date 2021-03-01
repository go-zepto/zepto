package authcore

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

type AuthResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type AuthResetPasswordResponse struct {
	Status bool    `json:"status"`
	Error  *string `json:"error,omitempty"`
}

type AuthRecoveryPasswordRequest struct {
	Email string `json:"email"`
}

type AuthRecoveryPasswordResponse struct {
	Status bool    `json:"status"`
	Error  *string `json:"error,omitempty"`
}

type AuthTokenResponse struct {
	Token *Token `json:"token"`
}

type AuthTokenErrorResponse struct {
	Error string `json:"error"`
}

type AuthDatasource interface {
	Auth(username string, password string) (PID, error)
	FindPIDByEmail(email string) (PID, error)
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

type AuthNotifier interface {
	NotifyResetPasswordToken(email string, token *Token, pid PID) error
	NotifyPasswordReseted(email string, pid PID) error
}
