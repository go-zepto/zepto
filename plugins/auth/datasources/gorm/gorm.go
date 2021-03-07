package gorm

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-zepto/zepto/plugins/auth/authcore"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrNilDB = errors.New("[auth - gorm] db is nil")
var ErrNilUserModel = errors.New("[auth - gorm] user model is nil")
var ErrNilUserTableDoesNotExist = errors.New("[auth - gorm] user table does not exist")
var ErrInvalidIDField = errors.New("[auth - gorm] invalid id field")
var ErrInvalidUsernameField = errors.New("[auth - gorm] invalid username field")
var ErrInvalidPasswordHashField = errors.New("[auth - gorm] invalid password hash field")

type GormAuthDatasource struct {
	db                *gorm.DB
	userModel         interface{}
	idField           string
	usernameField     string
	emailField        string
	passwordHashField string
}

type GormAuthDatasourceOptions struct {
	DB                *gorm.DB
	UserModel         interface{}
	IDField           string
	UsernameField     string
	EmailField        string
	PasswordHashField string
}

func getIDField(opts *GormAuthDatasourceOptions) string {
	if opts.IDField != "" {
		return opts.IDField
	}
	return "id"
}

func getUsernameField(opts *GormAuthDatasourceOptions) string {
	if opts.UsernameField != "" {
		return opts.UsernameField
	}
	return "username"
}

func getEmailField(opts *GormAuthDatasourceOptions) string {
	if opts.EmailField != "" {
		return opts.EmailField
	}
	return "email"
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
	if opts.DB == nil {
		panic(ErrNilDB)
	}
	gads := GormAuthDatasource{
		db:                opts.DB,
		userModel:         opts.UserModel,
		idField:           ensureGormFieldNaming(opts.DB, getIDField(&opts)),
		usernameField:     ensureGormFieldNaming(opts.DB, getUsernameField(&opts)),
		emailField:        ensureGormFieldNaming(opts.DB, getEmailField(&opts)),
		passwordHashField: ensureGormFieldNaming(opts.DB, getPasswordHashField(&opts)),
	}
	if err := gads.assertConditions(); err != nil {
		panic(err)
	}
	return &gads
}

func (gads *GormAuthDatasource) assertConditions() error {
	if gads.userModel == nil {
		return ErrNilUserModel
	}
	if !gads.db.Migrator().HasTable(gads.userModel) {
		return ErrNilUserTableDoesNotExist
	}
	if !gads.db.Migrator().HasColumn(gads.userModel, gads.idField) {
		return ErrInvalidIDField
	}
	if !gads.db.Migrator().HasColumn(gads.userModel, gads.usernameField) {
		return ErrInvalidUsernameField
	}
	if !gads.db.Migrator().HasColumn(gads.userModel, gads.passwordHashField) {
		return ErrInvalidPasswordHashField
	}
	return nil
}

func ensureInterfaceString(str interface{}) string {
	rv := reflect.ValueOf(str)
	if rv.Kind() == reflect.Ptr {
		return rv.Elem().String()
	}
	return rv.String()
}

func (gads *GormAuthDatasource) Auth(username string, password string) (authcore.PID, error) {
	if username == "" || password == "" {
		return nil, authcore.ErrUnauthorized
	}
	user := map[string]interface{}{}
	err := gads.db.Model(gads.userModel).Where(fmt.Sprintf("%s = ?", gads.usernameField), username).Find(&user).Error
	if err != nil {
		return nil, authcore.ErrUnauthorized
	}
	id, exists := user[gads.idField]
	if !exists {
		return nil, authcore.ErrUnauthorized
	}
	pwdHash, exists := user[gads.passwordHashField]
	pwdHashStr := ensureInterfaceString(pwdHash)
	if !exists || pwdHashStr == "" {
		return nil, authcore.ErrUnauthorized
	}
	if err := bcrypt.CompareHashAndPassword([]byte(pwdHashStr), []byte(password)); err != nil {
		return nil, authcore.ErrUnauthorized
	}
	return authcore.PID(id), nil
}

func (gads *GormAuthDatasource) ResetPassword(pid authcore.PID, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}
	passwordHashStr := string(passwordHash)
	return gads.db.Model(gads.userModel).Where(fmt.Sprintf("%s = ?", gads.idField), pid).Update(gads.passwordHashField, passwordHashStr).Error
}

func (gads *GormAuthDatasource) FindPIDByEmail(email string) (authcore.PID, error) {
	user := map[string]interface{}{}
	err := gads.db.Model(gads.userModel).Where(fmt.Sprintf("%s = ?", gads.emailField), email).Find(&user).Error
	if err != nil {
		return nil, authcore.ErrNotFound
	}
	id, exists := user[gads.idField]
	if !exists {
		return nil, authcore.ErrNotFound
	}
	return authcore.PID(id), nil
}
