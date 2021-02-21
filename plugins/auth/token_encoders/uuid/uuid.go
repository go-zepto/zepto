package uuid

import (
	"time"

	"github.com/go-zepto/zepto/plugins/auth"
	"github.com/google/uuid"
)

type UUIDTokenEncoder struct{}

func NewUUIDTokenEncoder() *UUIDTokenEncoder {
	return &UUIDTokenEncoder{}
}

// GenerateTokenFromPID generate a random unique token. The PID is not considered in this encoder
func (ute *UUIDTokenEncoder) GenerateTokenFromPID(pid auth.PID) (*auth.Token, error) {
	uuidv4, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	exp := now.Add(time.Hour * 24 * 30)
	token := auth.Token{
		Value:      uuidv4.String(),
		Expiration: &exp,
	}
	return &token, nil
}
