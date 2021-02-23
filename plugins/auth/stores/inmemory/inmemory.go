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
	authSessions             map[string]*Session
	passwordRecoverySessions map[string]*Session
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		authSessions:             make(map[string]*Session),
		passwordRecoverySessions: make(map[string]*Session),
	}
}

func (ims *InMemoryStore) storeToken(sessionsMap map[string]*Session, token *auth.Token, pid auth.PID) error {
	if token == nil || token.Value == "" || token.Expiration == nil {
		return auth.ErrInvalidToken
	}
	if pid == nil {
		return auth.ErrInvalidPID
	}
	sessionsMap[token.Value] = &Session{
		Token: token,
		PID:   pid,
	}
	return nil
}

func (ims *InMemoryStore) getTokenPID(sessionsMap map[string]*Session, tokenValue string) (auth.PID, error) {
	session, exists := sessionsMap[tokenValue]
	if !exists || session.Token == nil {
		return nil, auth.ErrUnauthorized
	}
	exp := *session.Token.Expiration
	if session.Token.Expiration == nil || exp.Before(time.Now()) {
		delete(sessionsMap, tokenValue)
		return nil, auth.ErrUnauthorized
	}
	return session.PID, nil
}

func (ims *InMemoryStore) StoreAuthToken(token *auth.Token, pid auth.PID) error {
	return ims.storeToken(ims.authSessions, token, pid)
}

func (ims *InMemoryStore) GetAuthTokenPID(tokenValue string) (auth.PID, error) {
	return ims.getTokenPID(ims.authSessions, tokenValue)
}

func (ims *InMemoryStore) StoreResetPasswordToken(token *auth.Token, pid auth.PID) error {
	panic("not implemented")
}

func (ims *InMemoryStore) GetResetPasswordTokenPID(token string) (auth.PID, error) {
	panic("not implemented")
}
