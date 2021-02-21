package uuid

import (
	"time"

	"github.com/go-zepto/zepto/plugins/auth"
	uuid "github.com/satori/go.uuid"
)

type UUIDTokenEncoder struct{}

func NewUUIDTokenEncoder() *UUIDTokenEncoder {
	return &UUIDTokenEncoder{}
}

func (ute *UUIDTokenEncoder) GenerateTokenFromPID(pid auth.PID) (*auth.Token, error) {
	u := uuid.NewV4()
	tokenVal := u.String()
	now := time.Now()
	exp := now.Add(time.Hour * 24 * 30)
	token := auth.Token{
		Value:      tokenVal,
		Expiration: &exp,
	}
	return &token, nil
}
