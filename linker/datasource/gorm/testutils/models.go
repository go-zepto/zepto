package testutils

import (
	"database/sql"
	"time"
)

type Person struct {
	ID          uint
	Name        string
	Email       *string
	Age         uint8
	Birthday    *time.Time
	Active      bool
	ActivatedAt sql.NullTime
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
