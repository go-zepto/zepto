package repository

import (
	"context"
	"testing"

	gormds "github.com/go-zepto/zepto/linker/datasource/gorm"
	"github.com/go-zepto/zepto/linker/datasource/gorm/testutils"
	"github.com/go-zepto/zepto/linker/filter"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func SetupRepository(db *gorm.DB) *Repository {
	return NewRepository(RepositoryConfig{
		Datasource: gormds.NewGormDatasource(db, &testutils.Person{}),
	})
}

func TestFindById(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepository(db)
	res, err := r.FindById(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	var p testutils.Person
	res.Decode(&p)
	assert.Equal(t, "Carlos Strand", p.Name)
}

func TestFindById_NotFound(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepository(db)
	res, err := r.FindById(context.Background(), 99)
	assert.EqualError(t, err, "record not found")
	assert.Nil(t, res)
}

func TestFindOne(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepository(db)
	where := map[string]interface{}{
		"age": map[string]interface{}{
			"eq": 65,
		},
	}
	res, err := r.FindOne(context.Background(), &filter.Filter{Where: &where})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	var p testutils.Person
	res.Decode(&p)
	assert.Equal(t, "Bill Gates", p.Name)
}

func TestFind(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepository(db)
	where := map[string]interface{}{
		"age": map[string]interface{}{
			"in": []uint{27, 65},
		},
	}
	res, err := r.Find(context.Background(), &filter.Filter{Where: &where})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(2), res.Count)
	var p []testutils.Person
	res.Decode(&p)
	assert.Equal(t, "Carlos Strand", p[0].Name)
	assert.Equal(t, "Bill Gates", p[1].Name)
}

func TestUpdateById(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepository(db)
	data := map[string]interface{}{
		"name": "Kal-el",
	}
	res, err := r.UpdateById(context.Background(), 3, data)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	var p testutils.Person
	res.Decode(&p)
	assert.Equal(t, "Kal-el", p.Name)
}

func TestUpdate(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepository(db)
	where := map[string]interface{}{
		"age": map[string]interface{}{
			"in": []uint{24, 27},
		},
	}
	data := map[string]interface{}{
		"name": "Young Person",
	}
	res, err := r.Update(context.Background(), &filter.Filter{Where: &where}, data)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(2), res.TotalAffected)
	for _, id := range []uint{1, 3} {
		res, err := r.FindById(context.Background(), id)
		assert.NoError(t, err)
		var p testutils.Person
		res.Decode(&p)
		assert.Equal(t, "Young Person", p.Name)
	}
}

func TestDestroyById(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepository(db)
	err := r.DestroyById(context.Background(), 3)
	assert.NoError(t, err)
	_, err = r.FindById(context.Background(), 3)
	assert.EqualError(t, err, "record not found")
}

func TestDestroy(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepository(db)
	where := map[string]interface{}{
		"age": map[string]interface{}{
			"in": []uint{24, 27},
		},
	}
	res, err := r.Destroy(context.Background(), &filter.Filter{Where: &where})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(2), res.TotalAffected)
	for _, id := range []uint{1, 3} {
		_, err := r.FindById(context.Background(), id)
		assert.EqualError(t, err, "record not found")
	}
}
