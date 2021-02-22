package inmemory

import (
	"time"

	"github.com/go-zepto/zepto/plugins/auth"
)

type Session struct {
	Token *auth.Token
	PID   auth.PID
}

type InMemoryStore struct {
	sessions map[string]*Session // map token value -> Session
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		sessions: make(map[string]*Session),
	}
}

func (ims *InMemoryStore) StoreAuthToken(token *auth.Token, pid auth.PID) error {
	if token == nil || token.Value == "" || token.Expiration == nil {
		return auth.ErrInvalidToken
	}
	if pid == nil {
		return auth.ErrInvalidPID
	}
	ims.sessions[token.Value] = &Session{
		Token: token,
		PID:   pid,
	}
	return nil
}

func (ims *InMemoryStore) GetAuthTokenPID(tokenValue string) (auth.PID, error) {
	session, exists := ims.sessions[tokenValue]
	if !exists || session.Token == nil {
		return nil, auth.ErrUnauthorized
	}
	exp := *session.Token.Expiration
	if session.Token.Expiration == nil || exp.Before(time.Now()) {
		delete(ims.sessions, tokenValue)
		return nil, auth.ErrUnauthorized
	}
	return session.PID, nil
}

func (ims *InMemoryStore) StoreResetPasswordToken(token *auth.Token, pid auth.PID) error {
	panic("not implemented")
}

func (ims *InMemoryStore) GetResetPasswordTokenPID(token string) (auth.PID, error) {
	panic("not implemented")
}
