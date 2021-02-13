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

func TestList(t *testing.T) {
	db := testutils.SetupGorm()
	gds := NewGormDatasource(db, &testutils.Person{})
	res, err := gds.List(datasource.QueryContext{})
	assert.NoError(t, err)
	assert.Equal(t, int64(3), res.Count)
}

func TestListWithWhere(t *testing.T) {
	db := testutils.SetupGorm()
	gds := NewGormDatasource(db, &testutils.Person{})
	res, err := gds.List(datasource.QueryContext{
		Filter: filter.Filter{
			Where: map[string]interface{}{
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
