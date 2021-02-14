package testutils

import (
	"time"
)

type Person struct {
	ID        uint
	Name      string
	Email     *string
	Age       uint8
	Birthday  *time.Time
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
