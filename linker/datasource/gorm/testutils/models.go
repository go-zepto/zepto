package testutils

import (
	"time"
)

type Person struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	Email     *string    `json:"email"`
	Age       uint8      `json:"age"`
	Birthday  *time.Time `json:"birthday"`
	Active    bool       `json:"active"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
