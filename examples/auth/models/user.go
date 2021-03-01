package models

type User struct {
	Model
	Email        string `json:"email"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}
