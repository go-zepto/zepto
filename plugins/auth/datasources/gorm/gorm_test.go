package gorm

import (
	"testing"

	"github.com/go-zepto/zepto/plugins/auth/datasources/gorm/testutils"
	"github.com/go-zepto/zepto/plugins/auth/datasources/gorm/testutils/models"
	"github.com/stretchr/testify/assert"
)

func shouldPanic(t *testing.T, f func()) {
	defer func() { recover() }()
	f()
	t.Errorf("should have panicked")
}

func TestNewGormAuthDatasoruce(t *testing.T) {
	db := testutils.SetupDB()
	gad := NewGormAuthDatasoruce(GormAuthDatasourceOptions{
		DB:        db,
		UserModel: &models.User{},
	})
	assert.NotNil(t, gad)
}

func TestNewGormAuthDatasoruce_ErrNilDB(t *testing.T) {
	defer func() { recover() }()
	NewGormAuthDatasoruce(GormAuthDatasourceOptions{
		UserModel: &models.User{},
	})
	t.Errorf("did not panic")
}

func TestNewGormAuthDatasoruce_ErrNilUserModel(t *testing.T) {
	defer func() { recover() }()
	db := testutils.SetupDB()
	NewGormAuthDatasoruce(GormAuthDatasourceOptions{
		DB: db,
	})
	t.Errorf("did not panic")
}

func TestNewGormAuthDatasoruce_NotMigratedModel(t *testing.T) {
	defer func() { recover() }()
	db := testutils.SetupDB()
	db.Migrator().DropTable(&models.User{})
	NewGormAuthDatasoruce(GormAuthDatasourceOptions{
		DB:        db,
		UserModel: &models.User{},
	})
	t.Errorf("did not panic")
}

func TestNewGormAuthDatasoruce_ErrInvalidUsernameField(t *testing.T) {
	defer func() { recover() }()
	db := testutils.SetupDB()
	NewGormAuthDatasoruce(GormAuthDatasourceOptions{
		DB:        db,
		UserModel: &models.UserWithInvalidUsernameField{},
	})
	t.Errorf("did not panic")
}

func TestNewGormAuthDatasoruce_ErrInvalidPasswordHashField(t *testing.T) {
	defer func() { recover() }()
	db := testutils.SetupDB()
	NewGormAuthDatasoruce(GormAuthDatasourceOptions{
		DB:        db,
		UserModel: &models.UserWithInvalidPasswordHashField{},
	})
	t.Errorf("did not panic")
}

func TestNewGormAuthDatasoruce_CustomUsernameAndPassword(t *testing.T) {
	db := testutils.SetupDB()
	gad := NewGormAuthDatasoruce(GormAuthDatasourceOptions{
		DB:                db,
		UserModel:         &models.CustomUser{},
		UsernameField:     "email",
		PasswordHashField: "pwd_hash",
	})
	assert.NotNil(t, gad)
}
