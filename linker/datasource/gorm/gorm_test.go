package gorm

import (
	"testing"

	"github.com/go-zepto/zepto/linker/datasource"
	"github.com/go-zepto/zepto/linker/datasource/gorm/testutils"
	"github.com/go-zepto/zepto/linker/filter"
	"github.com/stretchr/testify/assert"
)

func SetupTestDatasource() datasource.Datasource {
	db := testutils.SetupGorm()
	gds := NewGormDatasource(db, &testutils.Person{})
	return gds
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

func TestFind_Limit(t *testing.T) {
	db := testutils.SetupGorm()
	gds := NewGormDatasource(db, &testutils.Person{})
	limit := 1
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
	limit := 1
	skip := 1
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
	res, err := gds.FindOne(datasource.QueryContext{
		// Filter: &filter.Filter{
		// 	Limit: &limit,
		// 	Skip:  &skip,
		// 	Where: &map[string]interface{}{
		// 		"age": map[string]interface{}{
		// 			"gt": 0,
		// 		},
		// 	},
		// },
	})
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

func TestCreate(t *testing.T) {
	db := testutils.SetupGorm()
	gds := NewGormDatasource(db, &testutils.Person{})
	res, err := gds.Create(datasource.QueryContext{}, map[string]interface{}{
		"name": "Mariah Medeiros",
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	user := *res
	assert.Equal(t, "Mariah Medeiros", user["name"])
	assert.Equal(t, uint(4), user["id"])
}
