package models

type User struct {
	Model
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}
