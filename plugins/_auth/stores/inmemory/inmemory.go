package inmemory

import (
	"time"

	"github.com/go-zepto/zepto/plugins/auth/authcore"
)

type Session struct {
	Token *authcore.Token
	PID   authcore.PID
}

type InMemoryStore struct {
	authSessions          map[string]*Session
	resetPasswordSessions map[string]*Session
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		authSessions:          make(map[string]*Session),
		resetPasswordSessions: make(map[string]*Session),
	}
}

func (ims *InMemoryStore) storeToken(sessionsMap map[string]*Session, token *authcore.Token, pid authcore.PID) error {
	if token == nil || token.Value == "" || token.Expiration == nil {
		return authcore.ErrInvalidToken
	}
	if pid == nil {
		return authcore.ErrInvalidPID
	}
	sessionsMap[token.Value] = &Session{
		Token: token,
		PID:   pid,
	}
	return nil
}

func (ims *InMemoryStore) getTokenPID(sessionsMap map[string]*Session, tokenValue string) (authcore.PID, error) {
	session, exists := sessionsMap[tokenValue]
	if !exists || session.Token == nil {
		return nil, authcore.ErrUnauthorized
	}
	exp := *session.Token.Expiration
	if session.Token.Expiration == nil || exp.Before(time.Now()) {
		delete(sessionsMap, tokenValue)
		return nil, authcore.ErrUnauthorized
	}
	return session.PID, nil
}

func (ims *InMemoryStore) StoreAuthToken(token *authcore.Token, pid authcore.PID) error {
	return ims.storeToken(ims.authSessions, token, pid)
}

func (ims *InMemoryStore) GetAuthTokenPID(tokenValue string) (authcore.PID, error) {
	return ims.getTokenPID(ims.authSessions, tokenValue)
}

func (ims *InMemoryStore) StoreResetPasswordToken(token *authcore.Token, pid authcore.PID) error {
	return ims.storeToken(ims.resetPasswordSessions, token, pid)
}

func (ims *InMemoryStore) GetResetPasswordTokenPID(tokenValue string) (authcore.PID, error) {
	return ims.getTokenPID(ims.resetPasswordSessions, tokenValue)
}
