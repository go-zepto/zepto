package where

import (
	"testing"

	"github.com/go-zepto/zepto/linker/datasource/gorm/testutils"
	"github.com/stretchr/testify/assert"
)

func TestGormQuery(t *testing.T) {
	db := testutils.SetupGorm()
	filterJson := `
		{
			"name": {
				"eq": "Bill Gates"
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []testutils.Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 1)
	assert.Equal(t, people[0].Name, "Bill Gates")
}

func TestGormQueryOr(t *testing.T) {
	db := testutils.SetupGorm()
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
	var people []testutils.Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 2)
	assert.Equal(t, people[0].Name, "Carlos Strand")
	assert.Equal(t, people[1].Name, "Clark Kent")
}

func TestGormQuery_Boolean(t *testing.T) {
	db := testutils.SetupGorm()
	filterJson := `
		{
			"active": {
				"eq": true
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []testutils.Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 1)
	assert.Equal(t, people[0].Name, "Clark Kent")
}

func TestGormQuery_Integer(t *testing.T) {
	db := testutils.SetupGorm()
	filterJson := `
		{
			"age": {
				"eq": 65
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []testutils.Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 1)
	assert.Equal(t, people[0].Name, "Bill Gates")
}

func TestGormQuery_GreaterThan(t *testing.T) {
	db := testutils.SetupGorm()
	filterJson := `
		{
			"age": {
				"gt": 64
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []testutils.Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 1)
	assert.Equal(t, people[0].Name, "Bill Gates")
}

func TestGormQuery_GreaterEqualThan(t *testing.T) {
	db := testutils.SetupGorm()
	filterJson := `
		{
			"age": {
				"gte": 65
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []testutils.Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 1)
	assert.Equal(t, people[0].Name, "Bill Gates")
}

func TestGormQuery_LessThen(t *testing.T) {
	db := testutils.SetupGorm()
	filterJson := `
		{
			"age": {
				"lt": 25
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []testutils.Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 1)
	assert.Equal(t, people[0].Name, "Clark Kent")
}

func TestGormQuery_LessThenEqual(t *testing.T) {
	db := testutils.SetupGorm()
	filterJson := `
		{
			"age": {
				"lte": 24
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []testutils.Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 1)
	assert.Equal(t, people[0].Name, "Clark Kent")
}

func TestGormQuery_Between(t *testing.T) {
	db := testutils.SetupGorm()
	filterJson := `
		{
			"age": {
				"between": [20, 30]
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []testutils.Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 2)
	assert.Equal(t, people[0].Name, "Carlos Strand")
	assert.Equal(t, people[1].Name, "Clark Kent")
}

func TestGormQuery_BetweenDates1(t *testing.T) {
	db := testutils.SetupGorm()
	filterJson := `
		{
			"birthday": {
				"between": ["1992-02-10", "1994-02-10"]
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []testutils.Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 1)
	assert.Equal(t, people[0].Name, "Carlos Strand")
}

func TestGormQuery_BetweenDates2(t *testing.T) {
	db := testutils.SetupGorm()
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
	var people []testutils.Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 2)
	assert.Equal(t, people[0].Name, "Carlos Strand")
	assert.Equal(t, people[1].Name, "Bill Gates")

}
func TestGormQuery_IN(t *testing.T) {
	db := testutils.SetupGorm()
	filterJson := `
		{
			"age": {
				"in": [24, 65]
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []testutils.Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 2)
	assert.Equal(t, people[0].Name, "Bill Gates")
	assert.Equal(t, people[1].Name, "Clark Kent")
}

func TestGormQuery_IN_Invalid(t *testing.T) {
	filterJson := `
		{
			"age": {
				"in": "invalid-value"
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	_, err := w.ToSQL()
	assert.EqualError(t, err, "IN operator must be an array")
}

func TestGormQuery_NIN(t *testing.T) {
	db := testutils.SetupGorm()
	filterJson := `
		{
			"age": {
				"nin": [24, 65]
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []testutils.Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 1)
	assert.Equal(t, people[0].Name, "Carlos Strand")
}

func TestGormQuery_NIN_Invalid(t *testing.T) {
	filterJson := `
		{
			"age": {
				"nin": "invalid-value"
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	_, err := w.ToSQL()
	assert.EqualError(t, err, "NIN operator must be an array")
}

func TestGormQuery_LIKE(t *testing.T) {
	db := testutils.SetupGorm()
	filterJson := `
		{
			"name": {
				"like": "%ent%"
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []testutils.Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 1)
	assert.Equal(t, people[0].Name, "Clark Kent")
}

func TestGormQuery_NLIKE(t *testing.T) {
	db := testutils.SetupGorm()
	filterJson := `
		{
			"name": {
				"nlike": "%ent%"
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	var people []testutils.Person
	q, err := w.ToSQL()
	assert.NoError(t, err)
	err = db.Where(q.Text, q.Vars...).Find(&people).Error
	assert.NoError(t, err)
	assert.Len(t, people, 2)
	assert.Equal(t, people[0].Name, "Carlos Strand")
	assert.Equal(t, people[1].Name, "Bill Gates")
}
