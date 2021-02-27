package uuid

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateTokenFromPID(t *testing.T) {
	enc := NewUUIDTokenEncoder()
	token, err := enc.GenerateTokenFromPID(100)
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotNil(t, token.Expiration)
	exp := *token.Expiration
	assert.True(t, exp.After(time.Now()), "expiration date should be greater than now")
}
