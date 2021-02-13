package where

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/thriftrw/ptr"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Person struct {
	ID          uint
	Name        string
	Email       *string
	Age         uint8
	Birthday    *time.Time
	Active      bool
	ActivatedAt sql.NullTime
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

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

func TestGormQuery(t *testing.T) {
	db := SetupGorm()
	filterJson := `
		{
			"name": {
				"eq": "Bill Gates"
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 1)
	assert.Equal(t, people[0].Name, "Bill Gates")
}

func TestGormQueryOr(t *testing.T) {
	db := SetupGorm()
	filterJson := `
		{
			"or": [
				{
					"name": {
						"eq": "Carlos Strand"
					}
				},
				{
					"name": {
						"eq": "Clark Kent"
					}
				}
			]
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 2)
	assert.Equal(t, people[0].Name, "Carlos Strand")
	assert.Equal(t, people[1].Name, "Clark Kent")
}

func TestGormQuery_Boolean(t *testing.T) {
	db := SetupGorm()
	filterJson := `
		{
			"active": {
				"eq": true
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 1)
	assert.Equal(t, people[0].Name, "Clark Kent")
}

func TestGormQuery_Integer(t *testing.T) {
	db := SetupGorm()
	filterJson := `
		{
			"age": {
				"eq": 65
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 1)
	assert.Equal(t, people[0].Name, "Bill Gates")
}

func TestGormQuery_GreaterThan(t *testing.T) {
	db := SetupGorm()
	filterJson := `
		{
			"age": {
				"gt": 64
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 1)
	assert.Equal(t, people[0].Name, "Bill Gates")
}

func TestGormQuery_GreaterEqualThan(t *testing.T) {
	db := SetupGorm()
	filterJson := `
		{
			"age": {
				"gte": 65
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 1)
	assert.Equal(t, people[0].Name, "Bill Gates")
}

func TestGormQuery_LessThen(t *testing.T) {
	db := SetupGorm()
	filterJson := `
		{
			"age": {
				"lt": 25
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 1)
	assert.Equal(t, people[0].Name, "Clark Kent")
}

func TestGormQuery_LessThenEqual(t *testing.T) {
	db := SetupGorm()
	filterJson := `
		{
			"age": {
				"lte": 24
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 1)
	assert.Equal(t, people[0].Name, "Clark Kent")
}

func TestGormQuery_Between(t *testing.T) {
	db := SetupGorm()
	filterJson := `
		{
			"age": {
				"between": [20, 30]
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 2)
	assert.Equal(t, people[0].Name, "Carlos Strand")
	assert.Equal(t, people[1].Name, "Clark Kent")
}

func TestGormQuery_BetweenDates1(t *testing.T) {
	db := SetupGorm()
	filterJson := `
		{
			"birthday": {
				"between": ["1992-02-10", "1994-02-10"]
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 1)
	assert.Equal(t, people[0].Name, "Carlos Strand")
}

func TestGormQuery_BetweenDates2(t *testing.T) {
	db := SetupGorm()
	filterJson := `
		{
			"or": [
				{
					"birthday": {
						"between": ["1950-01-01", "1960-01-01"]
					}
				},
				{
					"birthday": {
						"between": ["1992-02-10", "1994-02-10"]
					}
				}
			]
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 2)
	assert.Equal(t, people[0].Name, "Carlos Strand")
	assert.Equal(t, people[1].Name, "Bill Gates")
}
