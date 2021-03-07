package gorm

import (
	"testing"
	"time"

	"github.com/go-zepto/zepto/plugins/linker/datasource"
	"github.com/go-zepto/zepto/plugins/linker/datasource/gorm/testutils"
	"github.com/go-zepto/zepto/plugins/linker/filter"
	"github.com/go-zepto/zepto/plugins/linker/filter/include"
	"github.com/stretchr/testify/assert"
	"go.uber.org/thriftrw/ptr"
)

func SetupTestDatasource() datasource.Datasource {
	db := testutils.SetupGorm()
	gds := NewGormDatasource(db, &testutils.Person{})
	return gds
}

func TestFields(t *testing.T) {
	db := testutils.SetupGorm()
	gds := NewGormDatasource(db, &testutils.Person{})
	fields := gds.Fields()
	var assertType = func(name string, fieldType string) {
		assert.Equal(
			t,
			datasource.Field{
				Name:     name,
				Type:     fieldType,
				Nullable: true,
			},
			fields[name],
		)
	}
	// Model Base
	assertType("id", "integer")
	assertType("created_at", "datetime")
	assertType("updated_at", "datetime")
	assertType("deleted_at", "datetime")
	// Person
	assertType("name", "text")
	assertType("email", "varchar(60)")
	assertType("age", "integer")
	assertType("birthday", "datetime")
	assertType("active", "numeric")
	assertType("city_id", "integer")
}

func TestFind(t *testing.T) {
	db := testutils.SetupGorm()
	gds := NewGormDatasource(db, &testutils.Person{})
	res, err := gds.Find(datasource.QueryContext{})
	assert.NoError(t, err)
	assert.Equal(t, int64(3), res.Count)
}

func TestFind_Where(t *testing.T) {
	db := testutils.SetupGorm()
	gds := NewGormDatasource(db, &testutils.Person{})
	res, err := gds.Find(datasource.QueryContext{
		Filter: &filter.Filter{
			Where: &map[string]interface{}{
				"name": map[string]interface{}{
					"eq": "Carlos Strand",
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), res.Count)
	assert.Equal(t, res.Data[0]["name"], "Carlos Strand")
}

func TestFind_Where_Not_Allowed_Field(t *testing.T) {
	db := testutils.SetupGorm()
	gds := NewGormDatasource(db, &testutils.Person{})
	_, err := gds.Find(datasource.QueryContext{
		Filter: &filter.Filter{
			Where: &map[string]interface{}{
				"some_unknown": map[string]interface{}{
					"eq": "Carlos Strand",
				},
			},
		},
	})
	assert.EqualError(t, err, "some_unknown field is not allowed")
}

func TestFind_Where_Include_1(t *testing.T) {
	db := testutils.SetupGorm()
	gds := NewGormDatasource(db, &testutils.Person{})
	res, err := gds.Find(datasource.QueryContext{
		Filter: &filter.Filter{
			Include: []include.Include{
				{
					Relation: "City",
				},
				{
					Relation: "Orders",
				},
			},
			Where: &map[string]interface{}{
				"name": map[string]interface{}{
					"eq": "Carlos Strand",
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), res.Count)
	assert.Equal(t, res.Data[0]["name"], "Carlos Strand")
	assert.NotNil(t, res.Data[0]["city"])
	assert.Equal(t, "Salvador", res.Data[0]["city"].(*testutils.City).Name)
	user := res.Data[0]
	assert.Len(t, user["orders"], 2)
	orders := user["orders"].([]testutils.Order)
	assert.Equal(t, "Carlos's Order (1)", orders[0].Name)
	assert.Equal(t, "Carlos's Order (2)", orders[1].Name)
}

func TestFind_Where_Include_2(t *testing.T) {
	db := testutils.SetupGorm()
	gds := NewGormDatasource(db, &testutils.Person{})
	res, err := gds.Find(datasource.QueryContext{
		Filter: &filter.Filter{
			Include: []include.Include{
				{
					Relation: "City",
				},
				{
					Relation: "Orders",
					Where: &map[string]interface{}{
						"approved": map[string]interface{}{
							"eq": true,
						},
					},
				},
			},
			Where: &map[string]interface{}{
				"name": map[string]interface{}{
					"eq": "Carlos Strand",
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), res.Count)
	assert.Equal(t, res.Data[0]["name"], "Carlos Strand")
	assert.NotNil(t, res.Data[0]["city"])
	assert.Equal(t, "Salvador", res.Data[0]["city"].(*testutils.City).Name)
	user := res.Data[0]
	assert.Len(t, user["orders"], 1)
	orders := user["orders"].([]testutils.Order)
	assert.Equal(t, "Carlos's Order (1)", orders[0].Name)
}

func TestFind_Limit(t *testing.T) {
	db := testutils.SetupGorm()
	gds := NewGormDatasource(db, &testutils.Person{})
	limit := int64(1)
	res, err := gds.Find(datasource.QueryContext{
		Filter: &filter.Filter{
			Limit: &limit,
			Where: &map[string]interface{}{
				"age": map[string]interface{}{
					"gt": 0,
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(3), res.Count)
	assert.Len(t, res.Data, 1)
	assert.Equal(t, res.Data[0]["name"], "Carlos Strand")
}

func TestFind_Skip(t *testing.T) {
	db := testutils.SetupGorm()
	gds := NewGormDatasource(db, &testutils.Person{})
	limit := int64(1)
	skip := int64(1)
	res, err := gds.Find(datasource.QueryContext{
		Filter: &filter.Filter{
			Limit: &limit,
			Skip:  &skip,
			Where: &map[string]interface{}{
				"age": map[string]interface{}{
					"gt": 0,
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(3), res.Count)
	assert.Len(t, res.Data, 1)
	assert.Equal(t, res.Data[0]["name"], "Bill Gates")
}

func TestFindOne(t *testing.T) {
	db := testutils.SetupGorm()
	gds := NewGormDatasource(db, &testutils.Person{})
	res, err := gds.FindOne(datasource.QueryContext{})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	user := *res
	assert.Equal(t, user["name"], "Carlos Strand")
}

func TestFindOne_Where(t *testing.T) {
	db := testutils.SetupGorm()
	gds := NewGormDatasource(db, &testutils.Person{})
	res, err := gds.FindOne(datasource.QueryContext{
		Filter: &filter.Filter{
			Where: &map[string]interface{}{
				"age": map[string]interface{}{
					"gt": 50,
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	user := *res
	assert.Equal(t, user["name"], "Bill Gates")
}

func TestFindOne_Where_Include_1(t *testing.T) {
	db := testutils.SetupGorm()
	gds := NewGormDatasource(db, &testutils.Person{})
	res, err := gds.FindOne(datasource.QueryContext{
		Filter: &filter.Filter{
			Include: []include.Include{
				{
					Relation: "City",
				},
				{
					Relation: "Orders",
				},
			},
			Where: &map[string]interface{}{
				"age": map[string]interface{}{
					"eq": 27,
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	user := *res
	assert.Equal(t, user["name"], "Carlos Strand")
	assert.NotNil(t, user["city"])
	cityArg := user["city"]
	city := cityArg.(*testutils.City)
	assert.Equal(t, city.Name, "Salvador")
	assert.Len(t, user["orders"], 2)
	orders := user["orders"].([]testutils.Order)
	assert.Equal(t, "Carlos's Order (1)", orders[0].Name)
	assert.Equal(t, "Carlos's Order (2)", orders[1].Name)
}

func TestFindOne_Where_Include_2(t *testing.T) {
	db := testutils.SetupGorm()
	gds := NewGormDatasource(db, &testutils.Person{})
	res, err := gds.FindOne(datasource.QueryContext{
		Filter: &filter.Filter{
			Include: []include.Include{
				{
					Relation: "City",
				},
				{
					Relation: "Orders",
					Where: &map[string]interface{}{
						"approved": map[string]interface{}{
							"eq": true,
						},
					},
				},
			},
			Where: &map[string]interface{}{
				"age": map[string]interface{}{
					"eq": 27,
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	user := *res
	assert.Equal(t, user["name"], "Carlos Strand")
	assert.NotNil(t, user["city"])
	cityArg := user["city"]
	city := cityArg.(*testutils.City)
	assert.Equal(t, city.Name, "Salvador")
	assert.Len(t, user["orders"], 1)
	orders := user["orders"].([]testutils.Order)
	assert.Equal(t, "Carlos's Order (1)", orders[0].Name)
}

func TestCreate(t *testing.T) {
	db := testutils.SetupGorm()
	gds := NewGormDatasource(db, &testutils.Person{})
	newDate := func(year int, month time.Month, date int) *time.Time {
		d := time.Date(year, month, date, 0, 0, 0, 0, time.Local)
		return &d
	}
	res, err := gds.Create(datasource.QueryContext{}, map[string]interface{}{
		"name":     "Bruce Wayne",
		"email":    "bruce@test.com",
		"age":      21,
		"birthday": newDate(1999, 10, 10),
		"active":   true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	user := *res
	assert.Equal(t, uint(4), user["id"])
	assert.Equal(t, "Bruce Wayne", user["name"])
	assert.Equal(t, ptr.String("bruce@test.com"), user["email"])
	assert.Equal(t, uint8(21), user["age"])
	assert.Equal(t, newDate(1999, 10, 10), user["birthday"])
	assert.Equal(t, true, user["active"])
}

func TestUpdate(t *testing.T) {
	db := testutils.SetupGorm()
	gds := NewGormDatasource(db, &testutils.Person{})
	newDate := func(year int, month time.Month, date int) *time.Time {
		d := time.Date(year, month, date, 0, 0, 0, 0, time.Local)
		return &d
	}
	res, err := gds.Update(datasource.QueryContext{
		Filter: &filter.Filter{
			Where: &map[string]interface{}{
				"id": map[string]interface{}{
					"eq": 1,
				},
			},
		},
	}, map[string]interface{}{
		"name":     "Batman",
		"email":    "batman@test.com",
		"age":      25,
		"birthday": newDate(1995, 10, 10),
		"active":   true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(1), res.TotalAffected)
	p := testutils.Person{}
	err = db.First(&p, 1).Error
	assert.NoError(t, err)
	assert.Equal(t, "Batman", p.Name)
	assert.Equal(t, ptr.String("batman@test.com"), p.Email)
	assert.Equal(t, uint8(25), p.Age)
	assert.Equal(t, newDate(1995, 10, 10).Unix(), p.Birthday.Unix())
	assert.Equal(t, true, p.Active)
}

func TestUpdate_Many(t *testing.T) {
	db := testutils.SetupGorm()
	gds := NewGormDatasource(db, &testutils.Person{})
	res, err := gds.Update(datasource.QueryContext{
		Filter: &filter.Filter{
			Where: &map[string]interface{}{
				"age": map[string]interface{}{
					"lt": 30,
				},
			},
		},
	}, map[string]interface{}{
		"name": "Young Person",
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(2), res.TotalAffected)
	people := []testutils.Person{}
	err = db.Find(&people).Error
	assert.NoError(t, err)
	var asserts = []struct {
		Name string
	}{
		{
			Name: "Young Person",
		},
		{
			Name: "Bill Gates",
		},
		{
			Name: "Young Person",
		},
	}
	for idx, p := range asserts {
		assert.Equal(t, p.Name, people[idx].Name)
	}
}

func TestDestroy(t *testing.T) {
	db := testutils.SetupGorm()
	gds := NewGormDatasource(db, &testutils.Person{})
	res, err := gds.Destroy(datasource.QueryContext{
		Filter: &filter.Filter{
			Where: &map[string]interface{}{
				"id": map[string]interface{}{
					"eq": 1,
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(1), res.TotalAffected)
	p := testutils.Person{}
	err = db.First(&p, 1).Error
	assert.EqualError(t, err, "record not found")
}

func TestDestroyMany(t *testing.T) {
	db := testutils.SetupGorm()
	gds := NewGormDatasource(db, &testutils.Person{})
	res, err := gds.Destroy(datasource.QueryContext{
		Filter: &filter.Filter{
			Where: &map[string]interface{}{
				"id": map[string]interface{}{
					"in": []uint{1, 2},
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(2), res.TotalAffected)
	p := []testutils.Person{}
	err = db.Find(&p, []uint{1, 2}).Error
	assert.NoError(t, err)
	assert.Len(t, p, 0)
}
