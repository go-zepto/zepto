package repository

import (
	"context"
	"testing"

	gormds "github.com/go-zepto/zepto/linker/datasource/gorm"
	"github.com/go-zepto/zepto/linker/datasource/gorm/testutils"
	"github.com/go-zepto/zepto/linker/filter"
	"github.com/go-zepto/zepto/linker/hooks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"
)

func SetupRepositoryDecoder(db *gorm.DB, operationHooks hooks.OperationHooks) *RepositoryDecoder {
	repo := NewRepository(RepositoryConfig{
		Datasource:     gormds.NewGormDatasource(db, &testutils.Person{}),
		OperationHooks: operationHooks,
	})
	return &RepositoryDecoder{
		Repo: repo,
	}
}

func TestDecoderFind_SliceAsDest(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepositoryDecoder(db, nil)
	var dest []testutils.Person
	filter := &filter.Filter{
		Where: &map[string]interface{}{
			"age": map[string]interface{}{
				"in": []int{27, 65},
			},
		},
	}
	err := r.Find(context.Background(), filter, &dest)
	assert.NoError(t, err)
	assert.Len(t, dest, 2)
	assert.Equal(t, "Carlos Strand", dest[0].Name)
	assert.Equal(t, "Bill Gates", dest[1].Name)
}

func TestDecoderFind_ListObjectAsDest(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepositoryDecoder(db, nil)
	var dest ListResult
	filter := &filter.Filter{
		Where: &map[string]interface{}{
			"age": map[string]interface{}{
				"in": []int{27, 65},
			},
		},
	}
	err := r.Find(context.Background(), filter, &dest)
	assert.NoError(t, err)
	assert.Equal(t, dest.Count, int64(2))
	assert.Len(t, dest.Data, 2)
	var people []testutils.Person
	err = dest.Data.Decode(&people)
	assert.NoError(t, err)
	assert.Len(t, people, 2)
	assert.Equal(t, "Carlos Strand", people[0].Name)
	assert.Equal(t, "Bill Gates", people[1].Name)
}

func TestDecoderFindOne(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepositoryDecoder(db, nil)
	var dest testutils.Person
	filter := &filter.Filter{
		Where: &map[string]interface{}{
			"age": map[string]interface{}{
				"eq": 65,
			},
		},
	}
	err := r.FindOne(context.Background(), filter, &dest)
	assert.NoError(t, err)
	assert.Equal(t, "Bill Gates", dest.Name)
}

func TestDecoderFindById(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepositoryDecoder(db, nil)
	var dest testutils.Person
	err := r.FindById(context.Background(), 3, &dest)
	assert.NoError(t, err)
	assert.Equal(t, "Clark Kent", dest.Name)
}

func TestDecoderCreate(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepositoryDecoder(db, nil)
	var out testutils.Person
	in := testutils.Person{
		Name:  "Clark Kent",
		Email: ptr.String("clark@kent.com"),
	}
	err := r.Create(context.Background(), in, &out)
	assert.NoError(t, err)
	assert.Equal(t, in.Name, out.Name)
	assert.Equal(t, in.Email, out.Email)
	assert.Equal(t, uint(4), out.ID)
}

func TestDecoderUpdateById(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepositoryDecoder(db, nil)
	var out testutils.Person
	in := testutils.Person{
		Name: "Carlos Strand (Updated)",
	}
	err := r.UpdateById(context.Background(), 1, in, &out)
	assert.NoError(t, err)
	assert.Equal(t, in.Name, out.Name)
	assert.Equal(t, uint(1), out.ID)
}

func TestDecoderUpdate(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepositoryDecoder(db, nil)
	var out ManyAffectedResult
	in := testutils.Person{
		Name: "Bill (Updated)",
	}
	filter := &filter.Filter{
		Where: &map[string]interface{}{
			"age": map[string]interface{}{
				"eq": 65,
			},
		},
	}
	err := r.Update(context.Background(), filter, in, &out)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), out.TotalAffected)
}

func TestDecoderDestroy(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepositoryDecoder(db, nil)
	var out ManyAffectedResult
	filter := &filter.Filter{
		Where: &map[string]interface{}{
			"age": map[string]interface{}{
				"in": []int{27, 65},
			},
		},
	}
	err := r.Destroy(context.Background(), filter, &out)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), out.TotalAffected)
}

func TestDecoderDestroyById(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepositoryDecoder(db, nil)
	err := r.DestroyById(context.Background(), 1)
	assert.NoError(t, err)
}
