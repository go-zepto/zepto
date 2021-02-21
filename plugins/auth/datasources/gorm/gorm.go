package gorm

import (
	"errors"

	"github.com/go-zepto/zepto/plugins/auth"
	"gorm.io/gorm"
)

var ErrNilDB = errors.New("[auth - gorm] db is nil")
var ErrNilUserModel = errors.New("[auth - gorm] user model is nil")
var ErrNilUserTableDoesNotExist = errors.New("[auth - gorm] user table does not exist")
var ErrInvalidUsernameField = errors.New("[auth - gorm] invalid username field")
var ErrInvalidPasswordHashField = errors.New("[auth - gorm] invalid password hash field")

type GormAuthDatasource struct {
	db                *gorm.DB
	userModel         interface{}
	usernameField     string
	passwordHashField string
}

type GormAuthDatasourceOptions struct {
	DB                *gorm.DB
	UserModel         interface{}
	UsernameField     string
	PasswordHashField string
}

func getUsernameField(opts *GormAuthDatasourceOptions) string {
	if opts.UsernameField != "" {
		return opts.UsernameField
	}
	return "username"
}

func getPasswordHashField(opts *GormAuthDatasourceOptions) string {
	if opts.PasswordHashField != "" {
		return opts.PasswordHashField
	}
	return "password_hash"
}

func ensureGormFieldNaming(db *gorm.DB, field string) string {
	return db.NamingStrategy.ColumnName("", field)
}

func NewGormAuthDatasoruce(opts GormAuthDatasourceOptions) *GormAuthDatasource {
	gads := GormAuthDatasource{
		db:                opts.DB,
		userModel:         opts.UserModel,
		usernameField:     ensureGormFieldNaming(opts.DB, getUsernameField(&opts)),
		passwordHashField: ensureGormFieldNaming(opts.DB, getPasswordHashField(&opts)),
	}
	if err := gads.assertConditions(); err != nil {
		panic(err)
	}
	return &gads
}

func (gads *GormAuthDatasource) assertConditions() error {
	if gads.db == nil {
		return ErrNilDB
	}
	if gads.userModel == nil {
		return ErrNilUserModel
	}
	if !gads.db.Migrator().HasTable(gads.userModel) {
		return ErrNilUserTableDoesNotExist
	}
	if !gads.db.Migrator().HasColumn(gads.userModel, gads.usernameField) {
		return ErrInvalidUsernameField
	}
	if !gads.db.Migrator().HasColumn(gads.userModel, gads.passwordHashField) {
		return ErrInvalidPasswordHashField
	}
	return nil
}

func (gads *GormAuthDatasource) Auth(username string, password string) (auth.PID, error) {
	return auth.PID(100), nil
}
