package inmemory

import (
	"testing"
	"time"

	"github.com/go-zepto/zepto/plugins/auth/authcore"
	"github.com/stretchr/testify/assert"
)

func TestNewInMemory(t *testing.T) {
	im := NewInMemoryStore()
	assert.NotNil(t, im)
	assert.NotNil(t, im.authSessions)
}

func TestInMemory_storeToken_authSessions(t *testing.T) {
	im := NewInMemoryStore()
	exp := time.Now().Add(time.Hour * 24)
	token := &authcore.Token{
		Expiration: &exp,
		Value:      "50559387-6dfa-4282-9d9d-efc5e9af3e72",
	}
	err := im.storeToken(im.authSessions, token, 120)
	assert.NoError(t, err)
	session := im.authSessions[token.Value]
	assert.Equal(t, token, session.Token)
	assert.Equal(t, 120, session.PID)
}

func TestInMemory_storeToken_resetPasswordSessions(t *testing.T) {
	im := NewInMemoryStore()
	exp := time.Now().Add(time.Hour * 24)
	token := &authcore.Token{
		Expiration: &exp,
		Value:      "50559387-6dfa-4282-9d9d-efc5e9af3e72",
	}
	err := im.storeToken(im.resetPasswordSessions, token, 140)
	assert.NoError(t, err)
	session := im.resetPasswordSessions[token.Value]
	assert.Equal(t, token, session.Token)
	assert.Equal(t, 140, session.PID)
}

func TestStoreAuthToken(t *testing.T) {
	im := NewInMemoryStore()
	exp := time.Now().Add(time.Hour * 24)
	token := &authcore.Token{
		Expiration: &exp,
		Value:      "50559387-6dfa-4282-9d9d-efc5e9af3e72",
	}
	err := im.StoreAuthToken(token, 160)
	assert.NoError(t, err)
	session := im.authSessions[token.Value]
	assert.Equal(t, token, session.Token)
	assert.Equal(t, 160, session.PID)
}

func TestStoreAuthToken_Invalid_Value(t *testing.T) {
	im := NewInMemoryStore()
	exp := time.Now().Add(time.Hour * 24)
	token := &authcore.Token{
		Expiration: &exp,
		Value:      "",
	}
	err := im.StoreAuthToken(token, 150)
	assert.EqualError(t, err, authcore.ErrInvalidToken.Error())
}

func TestStoreAuthToken_Invalid_Expiration(t *testing.T) {
	im := NewInMemoryStore()
	token := &authcore.Token{
		Value:      "50559387-6dfa-4282-9d9d-efc5e9af3e72",
		Expiration: nil,
	}
	err := im.StoreAuthToken(token, 150)
	assert.EqualError(t, err, authcore.ErrInvalidToken.Error())
}

func TestStoreAuthToken_Invalid_PID(t *testing.T) {
	im := NewInMemoryStore()
	exp := time.Now().Add(time.Hour * 24)
	token := &authcore.Token{
		Expiration: &exp,
		Value:      "50559387-6dfa-4282-9d9d-efc5e9af3e72",
	}
	err := im.StoreAuthToken(token, nil)
	assert.EqualError(t, err, authcore.ErrInvalidPID.Error())
}

func TestGetAuthTokenPID(t *testing.T) {
	im := NewInMemoryStore()
	exp := time.Now().Add(time.Hour * 24)
	token := &authcore.Token{
		Expiration: &exp,
		Value:      "50559387-6dfa-4282-9d9d-efc5e9af3e72",
	}
	err := im.StoreAuthToken(token, 160)
	assert.NoError(t, err)
	pid, err := im.GetAuthTokenPID(token.Value)
	assert.NoError(t, err)
	assert.Equal(t, 160, pid)
}

func TestGetAuthTokenPID_TokenNotFound(t *testing.T) {
	im := NewInMemoryStore()
	exp := time.Now().Add(time.Hour * 24)
	token := &authcore.Token{
		Expiration: &exp,
		Value:      "50559387-6dfa-4282-9d9d-efc5e9af3e72",
	}
	err := im.StoreAuthToken(token, 160)
	assert.NoError(t, err)
	pid, err := im.GetAuthTokenPID("another-token")
	assert.EqualError(t, err, authcore.ErrUnauthorized.Error())
	assert.Equal(t, nil, pid)
}

func TestGetAuthTokenPID_ExpiredToken(t *testing.T) {
	im := NewInMemoryStore()
	exp := time.Now().Add(time.Hour * -24)
	token := &authcore.Token{
		Expiration: &exp,
		Value:      "50559387-6dfa-4282-9d9d-efc5e9af3e72",
	}
	err := im.StoreAuthToken(token, 160)
	assert.NoError(t, err)
	pid, err := im.GetAuthTokenPID(token.Value)
	assert.EqualError(t, err, authcore.ErrUnauthorized.Error())
	assert.Equal(t, nil, pid)
	assert.Nil(t, im.authSessions[token.Value], "Expired token should be removed from session")
}

func TestStoreResetPasswordToken(t *testing.T) {
	im := NewInMemoryStore()
	exp := time.Now().Add(time.Hour * 24)
	token := &authcore.Token{
		Expiration: &exp,
		Value:      "50559387-6dfa-4282-9d9d-efc5e9af3e72",
	}
	err := im.StoreResetPasswordToken(token, 160)
	assert.NoError(t, err)
	session := im.resetPasswordSessions[token.Value]
	assert.Equal(t, token, session.Token)
	assert.Equal(t, 160, session.PID)
}

func TestGetResetPasswordTokenPID(t *testing.T) {
	im := NewInMemoryStore()
	exp := time.Now().Add(time.Hour * 24)
	token := &authcore.Token{
		Expiration: &exp,
		Value:      "50559387-6dfa-4282-9d9d-efc5e9af3e72",
	}
	err := im.StoreResetPasswordToken(token, 160)
	assert.NoError(t, err)
	pid, err := im.GetResetPasswordTokenPID(token.Value)
	assert.NoError(t, err)
	assert.Equal(t, 160, pid)
}
