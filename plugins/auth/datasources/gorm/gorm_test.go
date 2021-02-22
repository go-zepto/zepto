package gorm

import (
	"testing"

	"github.com/go-zepto/zepto/plugins/auth"
	"github.com/go-zepto/zepto/plugins/auth/datasources/gorm/testutils"
	"github.com/go-zepto/zepto/plugins/auth/datasources/gorm/testutils/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func shouldPanic(t *testing.T, f func()) {
	defer func() { recover() }()
	f()
	t.Errorf("should have panicked")
}

func TestEnsureInterfaceString(t *testing.T) {
	str := "abc"
	res := ensureInterfaceString(str)
	assert.Equal(t, "abc", res)
}

func TestEnsureInterfaceString_Pointer(t *testing.T) {
	str := "abc"
	res := ensureInterfaceString(&str)
	assert.Equal(t, "abc", res)
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

func TestNewGormAuthDatasoruce_ErrInvalidIDField(t *testing.T) {
	defer func() { recover() }()
	db := testutils.SetupDB()
	NewGormAuthDatasoruce(GormAuthDatasourceOptions{
		DB:        db,
		UserModel: &models.UserWithInvalidIDField{},
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
		IDField:           "p_id",
		UsernameField:     "email",
		PasswordHashField: "pwd_hash",
	})
	assert.NotNil(t, gad)
}

func TestNewGormAuthDatasoruce_CustomIDField(t *testing.T) {
	db := testutils.SetupDB()
	gad := NewGormAuthDatasoruce(GormAuthDatasourceOptions{
		DB:                db,
		UserModel:         &models.CustomUser{},
		IDField:           "p_id",
		UsernameField:     "email",
		PasswordHashField: "pwd_hash",
	})
	assert.NotNil(t, gad)
}

func TestNewGormAuthDatasoruce_Auth(t *testing.T) {
	db := testutils.SetupDB()
	hash, err := bcrypt.GenerateFromPassword([]byte("iamsuperman123"), 10)
	assert.NoError(t, err)
	db.Create(&models.User{
		Username:     "clark.kent",
		PasswordHash: string(hash),
	})
	gad := NewGormAuthDatasoruce(GormAuthDatasourceOptions{
		DB:        db,
		UserModel: &models.User{},
	})
	assert.NotNil(t, gad)
	pid, err := gad.Auth("clark.kent", "iamsuperman123")
	assert.NoError(t, err)
	assert.Equal(t, uint(1), pid)
}

func TestNewGormAuthDatasoruce_Auth_InvalidTableAfterInit(t *testing.T) {
	db := testutils.SetupDB()
	gad := NewGormAuthDatasoruce(GormAuthDatasourceOptions{
		DB:        db,
		UserModel: &models.User{},
	})
	db.Migrator().DropTable(&models.User{})
	assert.NotNil(t, gad)
	pid, err := gad.Auth("clark.kent", "iamsuperman123")
	assert.EqualError(t, err, auth.ErrUnauthorized.Error())
	assert.Nil(t, pid)
}

func TestNewGormAuthDatasoruce_Auth_EmptyHashInDatabase(t *testing.T) {
	db := testutils.SetupDB()
	db.Create(&models.User{
		Username:     "clark.kent",
		PasswordHash: "",
	})
	gad := NewGormAuthDatasoruce(GormAuthDatasourceOptions{
		DB:        db,
		UserModel: &models.User{},
	})
	assert.NotNil(t, gad)
	pid, err := gad.Auth("clark.kent", "iamsuperman123")
	assert.EqualError(t, err, auth.ErrUnauthorized.Error())
	assert.Nil(t, pid)
}

func TestNewGormAuthDatasoruce_Auth_EmptyInputs(t *testing.T) {
	db := testutils.SetupDB()
	hash, err := bcrypt.GenerateFromPassword([]byte("iamsuperman123"), 10)
	assert.NoError(t, err)
	db.Create(&models.User{
		Username:     "clark.kent",
		PasswordHash: string(hash),
	})
	gad := NewGormAuthDatasoruce(GormAuthDatasourceOptions{
		DB:        db,
		UserModel: &models.User{},
	})
	assert.NotNil(t, gad)
	pid, err := gad.Auth("", "")
	assert.EqualError(t, err, auth.ErrUnauthorized.Error())
	assert.Nil(t, pid)
}

func TestNewGormAuthDatasoruce_Auth_WrongUsername(t *testing.T) {
	db := testutils.SetupDB()
	hash, err := bcrypt.GenerateFromPassword([]byte("iamsuperman123"), 10)
	assert.NoError(t, err)
	db.Create(&models.User{
		Username:     "clark.kent",
		PasswordHash: string(hash),
	})
	gad := NewGormAuthDatasoruce(GormAuthDatasourceOptions{
		DB:        db,
		UserModel: &models.User{},
	})
	assert.NotNil(t, gad)
	pid, err := gad.Auth("kal.el", "iamsuperman123")
	assert.EqualError(t, err, auth.ErrUnauthorized.Error())
	assert.Nil(t, pid)
}

func TestNewGormAuthDatasoruce_Auth_WrongPassword(t *testing.T) {
	db := testutils.SetupDB()
	hash, err := bcrypt.GenerateFromPassword([]byte("iamsuperman123"), 10)
	assert.NoError(t, err)
	db.Create(&models.User{
		Username:     "clark.kent",
		PasswordHash: string(hash),
	})
	gad := NewGormAuthDatasoruce(GormAuthDatasourceOptions{
		DB:        db,
		UserModel: &models.User{},
	})
	assert.NotNil(t, gad)
	pid, err := gad.Auth("clark.kent", "wrongPassword")
	assert.EqualError(t, err, auth.ErrUnauthorized.Error())
	assert.Nil(t, pid)
}

func setupPointerFieldGad(t *testing.T) *GormAuthDatasource {
	db := testutils.SetupDB()
	hash, err := bcrypt.GenerateFromPassword([]byte("iamsuperman123"), 10)
	assert.NoError(t, err)
	username := "clark.kent"
	pwd_hash := string(hash)
	db.Create(&models.UserPointerFields{
		Username:     &username,
		PasswordHash: &pwd_hash,
	})
	gad := NewGormAuthDatasoruce(GormAuthDatasourceOptions{
		DB:        db,
		UserModel: &models.UserPointerFields{},
	})
	assert.NotNil(t, gad)
	return gad
}

func TestNewGormAuthDatasoruce_Auth_PointerFields(t *testing.T) {
	gad := setupPointerFieldGad(t)
	pid, err := gad.Auth("clark.kent", "iamsuperman123")
	assert.NoError(t, err)
	assert.Equal(t, uint(1), pid)
}

func TestNewGormAuthDatasoruce_Auth_PointerFields_WrongUsername(t *testing.T) {
	gad := setupPointerFieldGad(t)
	pid, err := gad.Auth("kal.el", "iamsuperman123")
	assert.EqualError(t, err, auth.ErrUnauthorized.Error())
	assert.Nil(t, pid)
}

func TestNewGormAuthDatasoruce_Auth_PointerFields_WrongPassword(t *testing.T) {
	gad := setupPointerFieldGad(t)
	pid, err := gad.Auth("clark.kent", "wrongPass")
	assert.EqualError(t, err, auth.ErrUnauthorized.Error())
	assert.Nil(t, pid)
}

func setupCustomFieldGad(t *testing.T) *GormAuthDatasource {
	db := testutils.SetupDB()
	hash, err := bcrypt.GenerateFromPassword([]byte("iamsuperman123"), 10)
	assert.NoError(t, err)
	err = db.Create(&models.CustomUser{
		PID:     200,
		Email:   "clark.kent@dc.com",
		PwdHash: string(hash),
	}).Error
	assert.NoError(t, err)
	gad := NewGormAuthDatasoruce(GormAuthDatasourceOptions{
		DB:                db,
		UserModel:         &models.CustomUser{},
		IDField:           "p_id",
		UsernameField:     "email",
		PasswordHashField: "pwd_hash",
	})
	assert.NotNil(t, gad)
	return gad
}

func TestNewGormAuthDatasoruce_Auth_CustomFields(t *testing.T) {
	gad := setupCustomFieldGad(t)
	pid, err := gad.Auth("clark.kent@dc.com", "iamsuperman123")
	assert.NoError(t, err)
	assert.Equal(t, int64(200), pid)
}

func TestNewGormAuthDatasoruce_Auth_CustomFields_WrongUsername(t *testing.T) {
	gad := setupCustomFieldGad(t)
	pid, err := gad.Auth("kal.el@dc.com", "iamsuperman123")
	assert.EqualError(t, err, auth.ErrUnauthorized.Error())
	assert.Nil(t, pid)
}

func TestNewGormAuthDatasoruce_Auth_CustomFields_WrongPassword(t *testing.T) {
	gad := setupCustomFieldGad(t)
	pid, err := gad.Auth("clark.kent@dc.com", "wrongPassword")
	assert.EqualError(t, err, auth.ErrUnauthorized.Error())
	assert.Nil(t, pid)
}
