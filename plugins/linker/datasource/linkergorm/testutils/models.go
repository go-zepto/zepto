package testutils

import (
	"time"
)

type Model struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type City struct {
	Model
	Name string `json:"name"`
}

type Order struct {
	Model
	Name                  string     `json:"name"`
	EstimatedShippingDate *time.Time `json:"estimated_shipping_date"`
	AmountInCents         int64      `json:"amount_in_cents"`
	Approved              bool       `json:"approved"`
	PersonID              uint       `json:"person_id"`
	Person                Person     `json:"person"`
}

type Person struct {
	Model
	Name     string     `json:"name"`
	Email    *string    `json:"email" gorm:"type:varchar(60)"`
	Age      uint8      `json:"age"`
	Birthday *time.Time `json:"birthday"`
	Active   bool       `json:"active"`
	CityID   uint       `json:"city_id"`
	City     *City      `json:"city"`
	Orders   []Order    `json:"orders"`
}
