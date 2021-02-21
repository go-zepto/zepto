package testutils

import (
	"log"
	"os"
	"time"

	"github.com/go-zepto/zepto/plugins/auth/datasources/gorm/testutils/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewTestDB() *gorm.DB {
	l := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Silent,
			Colorful:      true,
		},
	)
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: l,
	})
	if err != nil {
		panic(err)
	}
	return db
}

func SetupDB() *gorm.DB {
	db := NewTestDB()
	db.AutoMigrate(
		&models.User{},
		&models.CustomUser{},
	)
	return db
}
