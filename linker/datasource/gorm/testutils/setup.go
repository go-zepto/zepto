package testutils

import (
	"log"
	"os"
	"time"

	"go.uber.org/thriftrw/ptr"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupGorm() *gorm.DB {
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
	db.Migrator().DropTable(&Person{})
	db.AutoMigrate(
		&Person{},
	)
	newDate := func(year int, month time.Month, date int) *time.Time {
		d := time.Date(year, month, date, 0, 0, 0, 0, time.Local)
		return &d
	}
	persons := []Person{
		{
			Name:     "Carlos Strand",
			Email:    ptr.String("carlos@test.com"),
			Age:      27,
			Birthday: newDate(1993, 02, 10),
		},
		{
			Name:     "Bill Gates",
			Email:    ptr.String("bill@test.com"),
			Age:      65,
			Birthday: newDate(1955, 10, 28),
		},
		{
			Name:     "Clark Kent",
			Email:    ptr.String("clark@test.com"),
			Age:      24,
			Birthday: newDate(1996, 05, 12),
			Active:   true,
		},
	}
	db.Create(&persons)
	return db
}