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

func TestStoreToken(t *testing.T) {
	im := NewInMemoryStore()
	exp := time.Now().Add(time.Hour * 24)
	token := &auth.Token{
		Expiration: &exp,
		Value:      "50559387-6dfa-4282-9d9d-efc5e9af3e72",
	}
	err := im.StoreToken(token, 160)
	assert.NoError(t, err)
	session := im.sessions[token.Value]
	assert.Equal(t, token, session.Token)
	assert.Equal(t, 160, session.PID)
}

func TestStoreToken_Invalid_Value(t *testing.T) {
	im := NewInMemoryStore()
	exp := time.Now().Add(time.Hour * 24)
	token := &auth.Token{
		Expiration: &exp,
		Value:      "",
	}
	err := im.StoreToken(token, 150)
	assert.EqualError(t, err, auth.ErrInvalidToken.Error())
}

func TestStoreToken_Invalid_Expiration(t *testing.T) {
	im := NewInMemoryStore()
	token := &auth.Token{
		Value:      "50559387-6dfa-4282-9d9d-efc5e9af3e72",
		Expiration: nil,
	}
	err := im.StoreToken(token, 150)
	assert.EqualError(t, err, auth.ErrInvalidToken.Error())
}

func TestStoreToken_Invalid_PID(t *testing.T) {
	im := NewInMemoryStore()
	exp := time.Now().Add(time.Hour * 24)
	token := &auth.Token{
		Expiration: &exp,
		Value:      "50559387-6dfa-4282-9d9d-efc5e9af3e72",
	}
	err := im.StoreToken(token, nil)
	assert.EqualError(t, err, auth.ErrInvalidPID.Error())
}

func TestGetPIDFromTokenValue(t *testing.T) {
	im := NewInMemoryStore()
	exp := time.Now().Add(time.Hour * 24)
	token := &auth.Token{
		Expiration: &exp,
		Value:      "50559387-6dfa-4282-9d9d-efc5e9af3e72",
	}
	err := im.StoreToken(token, 160)
	assert.NoError(t, err)
	pid, err := im.GetPIDFromTokenValue(token.Value)
	assert.NoError(t, err)
	assert.Equal(t, 160, pid)
}

func TestGetPIDFromTokenValue_TokenNotFound(t *testing.T) {
	im := NewInMemoryStore()
	exp := time.Now().Add(time.Hour * 24)
	token := &auth.Token{
		Expiration: &exp,
		Value:      "50559387-6dfa-4282-9d9d-efc5e9af3e72",
	}
	err := im.StoreToken(token, 160)
	assert.NoError(t, err)
	pid, err := im.GetPIDFromTokenValue("another-token")
	assert.EqualError(t, err, auth.ErrUnauthorized.Error())
	assert.Equal(t, nil, pid)
}

func TestGetPIDFromTokenValue_ExpiredToken(t *testing.T) {
	im := NewInMemoryStore()
	exp := time.Now().Add(time.Hour * -24)
	token := &auth.Token{
		Expiration: &exp,
		Value:      "50559387-6dfa-4282-9d9d-efc5e9af3e72",
	}
	err := im.StoreToken(token, 160)
	assert.NoError(t, err)
	pid, err := im.GetPIDFromTokenValue(token.Value)
	assert.EqualError(t, err, auth.ErrUnauthorized.Error())
	assert.Equal(t, nil, pid)
	assert.Nil(t, im.sessions[token.Value], "Expired token should be removed from session")
}
