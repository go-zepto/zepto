package models

import (
	"time"
)

// gorm.Model base structure
type Model struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
