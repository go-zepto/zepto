package stores

import (
	"errors"
	"fmt"
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

func (ims *InMemoryStore) StoreToken(token *auth.Token, pid auth.PID) error {
	ims.sessions[token.Value] = &Session{
		Token: token,
		PID:   pid,
	}
	return nil
}

func (ims *InMemoryStore) GetPIDFromTokenValue(tokenValue string) (auth.PID, error) {
	session, exists := ims.sessions[tokenValue]
	if !exists || session.Token == nil {
		return errors.New("unauthorized"), nil
	}
	fmt.Println(session)
	exp := *session.Token.Expiration
	if session.Token.Expiration == nil || exp.Before(time.Now()) {
		return errors.New("unauthorized"), nil
	}
	return session.PID, nil
}
