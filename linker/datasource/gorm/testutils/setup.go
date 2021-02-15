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
	db.Migrator().DropTable(
		&Order{},
		&Person{},
		&City{},
	)
	db.AutoMigrate(
		&Order{},
		&City{},
		&Person{},
	)
	newDate := func(year int, month time.Month, date int) *time.Time {
		d := time.Date(year, month, date, 0, 0, 0, 0, time.Local)
		return &d
	}
	cities := []City{
		{
			Name: "Salvador",
		},
		{
			Name: "Seattle",
		},
		{
			Name: "Krypton",
		},
	}
	persons := []Person{
		{
			Name:     "Carlos Strand",
			Email:    ptr.String("carlos@test.com"),
			Age:      27,
			Birthday: newDate(1993, 02, 10),
			CityID:   1,
		},
		{
			Name:     "Bill Gates",
			Email:    ptr.String("bill@test.com"),
			Age:      65,
			Birthday: newDate(1955, 10, 28),
			CityID:   2,
		},
		{
			Name:     "Clark Kent",
			Email:    ptr.String("clark@test.com"),
			Age:      24,
			Birthday: newDate(1996, 05, 12),
			Active:   true,
			CityID:   3,
		},
	}
	orders := []Order{
		{
			Name:                  "Carlos's Order (1)",
			EstimatedShippingDate: newDate(2021, 2, 20),
			AmountInCents:         5000,
			Approved:              true,
			PersonID:              1,
		},
		{
			Name:                  "Carlos's Order (2)",
			EstimatedShippingDate: newDate(2021, 2, 22),
			AmountInCents:         12800,
			Approved:              false,
			PersonID:              1,
		},
		{
			Name:                  "Bill's Order (1)",
			EstimatedShippingDate: newDate(2021, 3, 10),
			AmountInCents:         600000,
			Approved:              true,
			PersonID:              2,
		},
		{
			Name:                  "Bill's Order (2)",
			EstimatedShippingDate: newDate(2021, 4, 20),
			AmountInCents:         900000,
			Approved:              true,
			PersonID:              2,
		},
	}
	db.Create(&cities)
	db.Create(&persons)
	db.Create(&orders)
	return db
}
