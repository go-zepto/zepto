package inmemory

import (
	"testing"
	"time"

	"github.com/go-zepto/zepto/plugins/auth"
	"github.com/stretchr/testify/assert"
)

func TestNewInMemory(t *testing.T) {
	im := NewInMemoryStore()
	assert.NotNil(t, im)
	assert.NotNil(t, im.sessions)
}

func TestStoreAuthToken(t *testing.T) {
	im := NewInMemoryStore()
	exp := time.Now().Add(time.Hour * 24)
	token := &auth.Token{
		Expiration: &exp,
		Value:      "50559387-6dfa-4282-9d9d-efc5e9af3e72",
	}
	err := im.StoreAuthToken(token, 160)
	assert.NoError(t, err)
	session := im.sessions[token.Value]
	assert.Equal(t, token, session.Token)
	assert.Equal(t, 160, session.PID)
}

func TestStoreAuthToken_Invalid_Value(t *testing.T) {
	im := NewInMemoryStore()
	exp := time.Now().Add(time.Hour * 24)
	token := &auth.Token{
		Expiration: &exp,
		Value:      "",
	}
	err := im.StoreAuthToken(token, 150)
	assert.EqualError(t, err, auth.ErrInvalidToken.Error())
}

func TestStoreAuthToken_Invalid_Expiration(t *testing.T) {
	im := NewInMemoryStore()
	token := &auth.Token{
		Value:      "50559387-6dfa-4282-9d9d-efc5e9af3e72",
		Expiration: nil,
	}
	err := im.StoreAuthToken(token, 150)
	assert.EqualError(t, err, auth.ErrInvalidToken.Error())
}

func TestStoreAuthToken_Invalid_PID(t *testing.T) {
	im := NewInMemoryStore()
	exp := time.Now().Add(time.Hour * 24)
	token := &auth.Token{
		Expiration: &exp,
		Value:      "50559387-6dfa-4282-9d9d-efc5e9af3e72",
	}
	err := im.StoreAuthToken(token, nil)
	assert.EqualError(t, err, auth.ErrInvalidPID.Error())
}

func TestGetAuthTokenPID(t *testing.T) {
	im := NewInMemoryStore()
	exp := time.Now().Add(time.Hour * 24)
	token := &auth.Token{
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
	token := &auth.Token{
		Expiration: &exp,
		Value:      "50559387-6dfa-4282-9d9d-efc5e9af3e72",
	}
	err := im.StoreAuthToken(token, 160)
	assert.NoError(t, err)
	pid, err := im.GetAuthTokenPID("another-token")
	assert.EqualError(t, err, auth.ErrUnauthorized.Error())
	assert.Equal(t, nil, pid)
}

func TestGetAuthTokenPID_ExpiredToken(t *testing.T) {
	im := NewInMemoryStore()
	exp := time.Now().Add(time.Hour * -24)
	token := &auth.Token{
		Expiration: &exp,
		Value:      "50559387-6dfa-4282-9d9d-efc5e9af3e72",
	}
	err := im.StoreAuthToken(token, 160)
	assert.NoError(t, err)
	pid, err := im.GetAuthTokenPID(token.Value)
	assert.EqualError(t, err, auth.ErrUnauthorized.Error())
	assert.Equal(t, nil, pid)
	assert.Nil(t, im.sessions[token.Value], "Expired token should be removed from session")
}
