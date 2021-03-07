package repository

import (
	"context"
	"testing"

	gormds "github.com/go-zepto/zepto/plugins/linker/datasource/gorm"
	"github.com/go-zepto/zepto/plugins/linker/datasource/gorm/testutils"
	"github.com/go-zepto/zepto/plugins/linker/filter"
	"github.com/go-zepto/zepto/plugins/linker/hooks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"
)

func SetupRepository(db *gorm.DB, operationHooks hooks.OperationHooks) *Repository {
	return NewRepository(RepositoryConfig{
		Datasource:     gormds.NewGormDatasource(db, &testutils.Person{}),
		OperationHooks: operationHooks,
	})
}

type OperationHooksMock struct {
	beforeInfo      *hooks.OperationHooksInfo
	beforeInfoStack []*hooks.OperationHooksInfo
	afterInfo       *hooks.OperationHooksInfo
	afterInfoStack  []*hooks.OperationHooksInfo
}

func (h *OperationHooksMock) BeforeOperation(info hooks.OperationHooksInfo) error {
	h.beforeInfo = &info
	h.beforeInfoStack = append(h.beforeInfoStack, &info)
	return nil
}

func (h *OperationHooksMock) AfterOperation(info hooks.OperationHooksInfo) error {
	h.afterInfo = &info
	h.afterInfoStack = append(h.afterInfoStack, &info)
	return nil
}

func TestFindById(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepository(db, nil)
	res, err := r.FindById(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	var p testutils.Person
	res.Decode(&p)
	assert.Equal(t, "Carlos Strand", p.Name)
}

func TestFindById_OperationHooks(t *testing.T) {
	db := testutils.SetupGorm()
	hooks := OperationHooksMock{}
	r := SetupRepository(db, &hooks)
	res, err := r.FindById(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	var p testutils.Person
	res.Decode(&p)
	assert.Equal(t, "Carlos Strand", p.Name)
	// Before Operation
	binfo := hooks.beforeInfo
	assert.Equal(t, "FindOne", binfo.Operation)
	expectedWhere := &map[string]interface{}{
		"id": map[string]interface{}{
			"eq": 1,
		},
	}
	assert.Equal(t, expectedWhere, binfo.QueryContext.Filter.Where)
	// After Operation
	ainfo := hooks.afterInfo
	assert.Equal(t, "FindOne", ainfo.Operation)
	assert.Equal(t, expectedWhere, ainfo.QueryContext.Filter.Where)
	assert.Equal(t, "1", *ainfo.ID)
}

func TestFindById_NotFound(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepository(db, nil)
	res, err := r.FindById(context.Background(), 99)
	assert.EqualError(t, err, "record not found")
	assert.Nil(t, res)
}

func TestFindOne(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepository(db, nil)
	where := map[string]interface{}{
		"age": &map[string]interface{}{
			"eq": 65,
		},
	}
	filter := &filter.Filter{Where: &where}
	res, err := r.FindOne(context.Background(), filter)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	var p testutils.Person
	res.Decode(&p)
	assert.Equal(t, "Bill Gates", p.Name)
}
func TestFindOne_OperationHooks(t *testing.T) {
	db := testutils.SetupGorm()
	hooks := OperationHooksMock{}
	r := SetupRepository(db, &hooks)
	where := map[string]interface{}{
		"age": &map[string]interface{}{
			"eq": 65,
		},
	}
	filter := &filter.Filter{Where: &where}
	res, err := r.FindOne(context.Background(), filter)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	var p testutils.Person
	res.Decode(&p)
	assert.Equal(t, "Bill Gates", p.Name)
	// Before Operation
	binfo := hooks.beforeInfo
	assert.Equal(t, "FindOne", binfo.Operation)
	assert.Equal(t, filter, binfo.QueryContext.Filter)
	// After Operation
	ainfo := hooks.afterInfo
	assert.Equal(t, "FindOne", ainfo.Operation)
	assert.Equal(t, filter, ainfo.QueryContext.Filter)
	assert.Equal(t, "2", *ainfo.ID)
}

func TestFind_OperationHooks(t *testing.T) {
	db := testutils.SetupGorm()
	hooks := OperationHooksMock{}
	r := SetupRepository(db, &hooks)
	where := map[string]interface{}{
		"age": map[string]interface{}{
			"in": []uint{27, 65},
		},
	}
	filter := filter.Filter{Where: &where}
	res, err := r.Find(context.Background(), &filter)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(2), res.Count)
	var p []testutils.Person
	res.Decode(&p)
	assert.Equal(t, "Carlos Strand", p[0].Name)
	assert.Equal(t, "Bill Gates", p[1].Name)
	// Before Operation
	binfo := hooks.beforeInfo
	assert.Equal(t, "Find", binfo.Operation)
	assert.Equal(t, &filter, binfo.QueryContext.Filter)
	// After Operation
	ainfo := hooks.afterInfo
	assert.Equal(t, "Find", ainfo.Operation)
	assert.Equal(t, &filter, ainfo.QueryContext.Filter)
	assert.Nil(t, ainfo.ID)
}

func TestCreate(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepository(db, nil)
	data := map[string]interface{}{
		"name":  "Clark Kent",
		"email": ptr.String("clark@kent.com"),
	}
	res, err := r.Create(context.Background(), data)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	var d map[string]interface{}
	res.Decode(&d)
	assert.Equal(t, data["name"], d["name"])
	assert.Equal(t, data["email"], d["email"])
}

func TestCreateFromModel(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepository(db, nil)
	data := testutils.Person{
		Name:  "New Name",
		Email: ptr.String("email@test.com"),
	}
	res, err := r.Create(context.Background(), &data)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	var d testutils.Person
	res.Decode(&d)
	assert.Equal(t, d.ID, uint(4))
	assert.Equal(t, data.Name, d.Name)
	assert.Equal(t, data.Email, d.Email)
}

func TestCreate_OperationHooks(t *testing.T) {
	db := testutils.SetupGorm()
	hooks := OperationHooksMock{}
	r := SetupRepository(db, &hooks)
	data := map[string]interface{}{
		"name":  "Clark Kent",
		"email": ptr.String("clark@kent.com"),
	}
	res, err := r.Create(context.Background(), data)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	var d map[string]interface{}
	res.Decode(&d)
	assert.Equal(t, data["name"], d["name"])
	assert.Equal(t, data["email"], d["email"])
	// Before Operation
	binfo := hooks.beforeInfo
	assert.Equal(t, "Create", binfo.Operation)
	// After Operation
	ainfo := hooks.afterInfo
	assert.Equal(t, "Create", ainfo.Operation)
	assert.NotNil(t, ainfo.Data)
	dt := *ainfo.Data
	assert.Equal(t, dt["name"], data["name"])
	assert.Equal(t, dt["email"], data["email"])
	assert.Nil(t, ainfo.ID)
}

func TestUpdateById(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepository(db, nil)
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

func TestUpdateById_OperationHook(t *testing.T) {
	db := testutils.SetupGorm()
	hooks := OperationHooksMock{}
	r := SetupRepository(db, &hooks)
	data := map[string]interface{}{
		"name": "Kal-el",
	}
	res, err := r.UpdateById(context.Background(), 3, data)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	var p testutils.Person
	res.Decode(&p)
	assert.Equal(t, "Kal-el", p.Name)
	// Before Operation
	assert.Len(t, hooks.beforeInfoStack, 2)
	assert.Equal(t, "Update", hooks.beforeInfoStack[0].Operation)
	assert.Equal(t, "FindOne", hooks.beforeInfoStack[1].Operation)
	// After Operation
	assert.Len(t, hooks.afterInfoStack, 2)
	assert.Equal(t, "Update", hooks.afterInfoStack[0].Operation)
	assert.Equal(t, "FindOne", hooks.afterInfoStack[1].Operation)
	updateRes := hooks.afterInfoStack[0].Data
	assert.Equal(t, &map[string]interface{}{"total_affected": int64(1)}, updateRes)
}

func TestUpdate(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepository(db, nil)
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

func TestUpdateFromModel(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepository(db, nil)
	where := map[string]interface{}{
		"age": map[string]interface{}{
			"in": []uint{24, 27},
		},
	}
	data := testutils.Person{
		Name: "Young Person Directly from GORM Model",
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
		assert.Equal(t, "Young Person Directly from GORM Model", p.Name)
	}
}

func TestUpdate_OperationHook(t *testing.T) {
	db := testutils.SetupGorm()
	hooks := OperationHooksMock{}
	r := SetupRepository(db, &hooks)
	where := map[string]interface{}{
		"age": map[string]interface{}{
			"in": []uint{24, 27},
		},
	}
	data := map[string]interface{}{
		"name": "Young Person",
	}
	filter := &filter.Filter{Where: &where}
	res, err := r.Update(context.Background(), filter, data)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(2), res.TotalAffected)
	// Before Operation
	binfo := hooks.beforeInfo
	assert.Equal(t, "Update", binfo.Operation)
	assert.Equal(t, filter, binfo.QueryContext.Filter)
	// After Operation
	ainfo := hooks.afterInfo
	assert.Equal(t, "Update", ainfo.Operation)
	assert.Equal(t, filter, ainfo.QueryContext.Filter)
	assert.Equal(t, &map[string]interface{}{"total_affected": int64(2)}, ainfo.Data)
	assert.Nil(t, ainfo.ID)
}

func TestDestroyById(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepository(db, nil)
	err := r.DestroyById(context.Background(), 3)
	assert.NoError(t, err)
	_, err = r.FindById(context.Background(), 3)
	assert.EqualError(t, err, "record not found")
}

func TestDestroy(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepository(db, nil)
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

func TestDestroy_OperationHook(t *testing.T) {
	db := testutils.SetupGorm()
	hooks := OperationHooksMock{}
	r := SetupRepository(db, &hooks)
	where := map[string]interface{}{
		"age": map[string]interface{}{
			"in": []uint{24, 27},
		},
	}
	filter := &filter.Filter{Where: &where}
	res, err := r.Destroy(context.Background(), filter)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(2), res.TotalAffected)
	// Before Operation
	binfo := hooks.beforeInfo
	assert.Equal(t, "Destroy", binfo.Operation)
	assert.Equal(t, filter, binfo.QueryContext.Filter)
	// After Operation
	ainfo := hooks.afterInfo
	assert.Equal(t, "Destroy", ainfo.Operation)
	assert.Equal(t, filter, ainfo.QueryContext.Filter)
	assert.Equal(t, &map[string]interface{}{"total_affected": int64(2)}, ainfo.Data)
	assert.Nil(t, ainfo.ID)
}
