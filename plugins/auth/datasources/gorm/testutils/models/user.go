package models

type User struct {
	Model
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}

type UserWithInvalidUsernameField struct {
	Model
	PasswordHash string `json:"-"`
}

type UserWithInvalidPasswordHashField struct {
	Model
	Username string `json:"username"`
}

type CustomUser struct {
	Model
	Email   string `json:"email"`
	PwdHash string `json:"-"`
}
