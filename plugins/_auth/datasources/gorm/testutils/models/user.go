package models

type User struct {
	Model
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type UserWithInvalidIDField struct {
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
	PID     int64  `gorm:"primaryKey,autoIncrement" json:"p_id"`
	Email   string `json:"email"`
	PwdHash string `json:"-"`
}

type UserPointerFields struct {
	Model
	Username     *string `json:"username"`
	PasswordHash *string `json:"-"`
}
