package linker

import (
	"context"
	"testing"

	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/plugins/linker/datasource"
	gormds "github.com/go-zepto/zepto/plugins/linker/datasource/gorm"
	"github.com/go-zepto/zepto/plugins/linker/datasource/gorm/testutils"
	"github.com/go-zepto/zepto/plugins/linker/filter"
	"github.com/stretchr/testify/assert"
	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"
)

func SetupRepository(db *gorm.DB, operationHooks OperationHooks) *Repository {
	z := zepto.NewZepto()
	return NewRepository(RepositoryConfig{
		Linker:         NewLinker(z.Router("/linker/api")),
		ResourceName:   "MyResource",
		Datasource:     gormds.NewGormDatasource(db, &testutils.Person{}),
		OperationHooks: operationHooks,
	})
}

func assertOperationInfo(t *testing.T, expected *OperationHooksInfo, info *OperationHooksInfo) {
	assert.NotNil(t, info)
	assert.NotNil(t, info.Linker)
	if expected.QueryContext != nil {
		assert.Equal(t, expected.QueryContext.Filter.Where, info.QueryContext.Filter.Where)
	}
	assert.Equal(t, expected.Operation, info.Operation)
	assert.Equal(t, expected.Data, info.Data)
	assert.Equal(t, expected.ResourceID, info.ResourceID)
	assert.Equal(t, expected.ResourceName, info.ResourceName)
}

type OperationHooksMock struct {
	beforeInfo      *OperationHooksInfo
	beforeInfoStack []*OperationHooksInfo
	afterInfo       *OperationHooksInfo
	afterInfoStack  []*OperationHooksInfo
}

func (h *OperationHooksMock) BeforeOperation(info OperationHooksInfo) error {
	h.beforeInfo = &info
	h.beforeInfoStack = append(h.beforeInfoStack, &info)
	return nil
}

func (h *OperationHooksMock) AfterOperation(info OperationHooksInfo) error {
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
	assertOperationInfo(t, &OperationHooksInfo{
		Operation:    "FindOne",
		ResourceName: "MyResource",
		ResourceID:   ptr.String("1"),
		Data:         nil,
		QueryContext: &datasource.QueryContext{
			Filter: &filter.Filter{
				Where: &map[string]interface{}{
					"id": map[string]interface{}{
						"eq": 1,
					},
				},
			},
		},
	}, binfo)
	// After Operation
	ainfo := hooks.afterInfo
	assertOperationInfo(t, &OperationHooksInfo{
		Operation:    "FindOne",
		ResourceName: "MyResource",
		ResourceID:   ptr.String("1"),
		Data:         nil,
		QueryContext: &datasource.QueryContext{
			Filter: &filter.Filter{
				Where: &map[string]interface{}{
					"id": map[string]interface{}{
						"eq": 1,
					},
				},
			},
		},
	}, ainfo)
	assert.NotNil(t, ainfo.Result)
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
	f := &filter.Filter{Where: &where}
	res, err := r.FindOne(context.Background(), f)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	var p testutils.Person
	res.Decode(&p)
	assert.Equal(t, "Bill Gates", p.Name)
	// Before Operation
	binfo := hooks.beforeInfo
	assertOperationInfo(t, &OperationHooksInfo{
		Operation:    "FindOne",
		ResourceName: "MyResource",
		Data:         nil,
		QueryContext: &datasource.QueryContext{
			Filter: &filter.Filter{
				Where: &where,
			},
		},
	}, binfo)
	// After Operation
	ainfo := hooks.afterInfo
	assertOperationInfo(t, &OperationHooksInfo{
		Operation:    "FindOne",
		ResourceName: "MyResource",
		ResourceID:   ptr.String("2"),
		Data:         nil,
		QueryContext: &datasource.QueryContext{
			Filter: &filter.Filter{
				Where: &where,
			},
		},
	}, ainfo)
	assert.NotNil(t, ainfo.Result)
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
	f := filter.Filter{Where: &where}
	res, err := r.Find(context.Background(), &f)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(2), res.Count)
	var p []testutils.Person
	res.Decode(&p)
	assert.Equal(t, "Carlos Strand", p[0].Name)
	assert.Equal(t, "Bill Gates", p[1].Name)
	// Before Operation
	binfo := hooks.beforeInfo
	assertOperationInfo(t, &OperationHooksInfo{
		Operation:    "Find",
		ResourceName: "MyResource",
		Data:         nil,
		QueryContext: &datasource.QueryContext{
			Filter: &filter.Filter{
				Where: &where,
			},
		},
	}, binfo)
	// After Operation
	ainfo := hooks.afterInfo
	assertOperationInfo(t, &OperationHooksInfo{
		Operation:    "Find",
		ResourceName: "MyResource",
		ResourceID:   nil,
		Data:         nil,
		QueryContext: &datasource.QueryContext{
			Filter: &filter.Filter{
				Where: &where,
			},
		},
	}, ainfo)
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
	assertOperationInfo(t, &OperationHooksInfo{
		Operation:    "Create",
		ResourceName: "MyResource",
		Data:         &data,
	}, binfo)
	// After Operation
	ainfo := hooks.afterInfo
	assertOperationInfo(t, &OperationHooksInfo{
		Operation:    "Create",
		ResourceName: "MyResource",
		ResourceID:   ptr.String("4"),
		Data:         &data,
	}, ainfo)
	assert.NotNil(t, ainfo.Result)
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
	assert.Len(t, hooks.beforeInfoStack, 2)
	assert.Equal(t, "Update", hooks.beforeInfoStack[0].Operation)
	assert.Equal(t, "FindOne", hooks.beforeInfoStack[1].Operation)
	assert.Len(t, hooks.afterInfoStack, 2)
	assert.Equal(t, "Update", hooks.afterInfoStack[0].Operation)
	assert.Equal(t, "FindOne", hooks.afterInfoStack[1].Operation)
	// Before Operation
	binfo := hooks.beforeInfoStack[0]
	assertOperationInfo(t, &OperationHooksInfo{
		Operation:    "Update",
		ResourceName: "MyResource",
		ResourceID:   ptr.String("3"),
		Data:         &data,
	}, binfo)
	// After Operation
	ainfo := hooks.afterInfoStack[0]
	assertOperationInfo(t, &OperationHooksInfo{
		Operation:    "Update",
		ResourceName: "MyResource",
		ResourceID:   ptr.String("3"),
		Data:         &data,
	}, ainfo)
	assert.NotNil(t, ainfo.Result)
	updateRes := hooks.afterInfoStack[0].Result
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
	f := &filter.Filter{Where: &where}
	res, err := r.Update(context.Background(), f, data)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(2), res.TotalAffected)
	// Before Operation
	binfo := hooks.beforeInfoStack[0]
	assertOperationInfo(t, &OperationHooksInfo{
		Operation:    "Update",
		ResourceName: "MyResource",
		QueryContext: &datasource.QueryContext{
			Filter: &filter.Filter{Where: &where},
		},
		Data: &data,
	}, binfo)
	// After Operation
	ainfo := hooks.afterInfoStack[0]
	assertOperationInfo(t, &OperationHooksInfo{
		Operation:    "Update",
		ResourceName: "MyResource",
		QueryContext: &datasource.QueryContext{
			Filter: &filter.Filter{Where: &where},
		},
		ResourceID: nil,
		Data:       &data,
	}, ainfo)
	assert.NotNil(t, ainfo.Result)
}

func TestDestroyById(t *testing.T) {
	db := testutils.SetupGorm()
	r := SetupRepository(db, nil)
	err := r.DestroyById(context.Background(), 3)
	assert.NoError(t, err)
	_, err = r.FindById(context.Background(), 3)
	assert.EqualError(t, err, "record not found")
}

func TestDestroyById_OperationHook(t *testing.T) {
	db := testutils.SetupGorm()
	hooks := OperationHooksMock{}
	r := SetupRepository(db, &hooks)
	err := r.DestroyById(context.Background(), 2)
	assert.NoError(t, err)
	// Before Operation
	binfo := hooks.beforeInfoStack[0]
	assertOperationInfo(t, &OperationHooksInfo{
		Operation:    "Destroy",
		ResourceName: "MyResource",
		ResourceID:   ptr.String("2"),
	}, binfo)
	// After Operation
	ainfo := hooks.afterInfoStack[0]
	assertOperationInfo(t, &OperationHooksInfo{
		Operation:    "Destroy",
		ResourceName: "MyResource",
		ResourceID:   ptr.String("2"),
	}, ainfo)
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
	f := &filter.Filter{Where: &where}
	res, err := r.Destroy(context.Background(), f)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(2), res.TotalAffected)
	// Before Operation
	binfo := hooks.beforeInfoStack[0]
	assertOperationInfo(t, &OperationHooksInfo{
		Operation:    "Destroy",
		ResourceName: "MyResource",
		QueryContext: &datasource.QueryContext{
			Filter: &filter.Filter{Where: &where},
		},
	}, binfo)
	// After Operation
	ainfo := hooks.afterInfoStack[0]
	assertOperationInfo(t, &OperationHooksInfo{
		Operation:    "Destroy",
		ResourceName: "MyResource",
		QueryContext: &datasource.QueryContext{
			Filter: &filter.Filter{Where: &where},
		},
		ResourceID: nil,
	}, ainfo)
}
