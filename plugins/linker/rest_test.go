package linker

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-zepto/zepto"
	gormds "github.com/go-zepto/zepto/plugins/linker/datasource/gorm"
	"github.com/go-zepto/zepto/plugins/linker/datasource/gorm/testutils"
	"github.com/go-zepto/zepto/plugins/linker/utils"
	"github.com/go-zepto/zepto/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/thriftrw/ptr"
	"gorm.io/gorm"
)

type TestKit struct {
	app    *zepto.Zepto
	router *web.Router
	db     *gorm.DB
	linker *Linker
}

func NewTestKit(t *testing.T) TestKit {
	r := require.New(t)
	app := zepto.NewZepto()
	apiRouter := app.Router("/api")
	r.NotNil(apiRouter)
	db := testutils.SetupGorm()
	linker := NewLinker(apiRouter)
	return TestKit{
		app:    app,
		router: apiRouter,
		db:     db,
		linker: linker,
	}
}

func assertRemoteInfo(t *testing.T, expected *RemoteHooksInfo, info *RemoteHooksInfo) {
	assert.NotNil(t, info)
	assert.NotNil(t, info.Ctx)
	assert.NotNil(t, info.Linker)
	assert.Equal(t, expected.Endpoint, info.Endpoint)
	assert.Equal(t, expected.Data, info.Data)
	assert.Equal(t, expected.ResourceID, info.ResourceID)
	assert.Equal(t, expected.ResourceName, info.ResourceName)
}

func TestNewResource(t *testing.T) {
	r := require.New(t)
	k := NewTestKit(t)
	k.linker.AddResource(Resource{
		Name:       "Person",
		Datasource: gormds.NewGormDatasource(k.db, &testutils.Person{}),
	})
	r.Len(k.linker.repositories, 1)
	r.NotNil(k.linker.Repository("Person"))
	r.Nil(k.linker.Repository("Unknown"))
}

type RemoteHooksMock struct {
	beforeRemoteCallInfo *RemoteHooksInfo
	afterRemoteCallInfo  *RemoteHooksInfo
}

func (r *RemoteHooksMock) BeforeRemote(info RemoteHooksInfo) error {
	r.beforeRemoteCallInfo = &info
	return nil
}

func (r *RemoteHooksMock) AfterRemote(info RemoteHooksInfo) error {
	r.afterRemoteCallInfo = &info
	data := *info.Result
	data["custom_field"] = "perfect!"
	return nil
}

type RemoteHooksMockError struct{}

func (r *RemoteHooksMockError) BeforeRemote(info RemoteHooksInfo) error {
	info.Ctx.SetStatus(400)
	return errors.New("ups! strange error")
}

func (r *RemoteHooksMockError) AfterRemote(info RemoteHooksInfo) error {
	return nil
}

func TestBeforeRemoteHooksList(t *testing.T) {
	r := require.New(t)
	k := NewTestKit(t)
	h := RemoteHooksMock{}
	k.linker.AddResource(Resource{
		Name:        "Person",
		Datasource:  gormds.NewGormDatasource(k.db, &testutils.Person{}),
		RemoteHooks: &h,
	})
	w := httptest.NewRecorder()
	k.app.InitApp()
	k.app.ServeHTTP(w, httptest.NewRequest("GET", "/api/people", nil))
	type Res struct {
		Data        []testutils.Person `json:"data"`
		Count       int                `json:"count"`
		CustomField string             `json:"custom_field"`
	}
	var res Res
	json.Unmarshal(w.Body.Bytes(), &res)
	r.Equal(http.StatusOK, w.Code)
	r.Equal(3, res.Count)
	r.Equal("perfect!", res.CustomField)
	r.Equal("Carlos Strand", res.Data[0].Name)
	r.Equal("Bill Gates", res.Data[1].Name)
	r.Equal("Clark Kent", res.Data[2].Name)
	// Ensure BeforeRemoteHooks was called
	binfo := h.beforeRemoteCallInfo
	assertRemoteInfo(t, &RemoteHooksInfo{
		Endpoint:     "List",
		ResourceName: "Person",
		ResourceID:   nil,
		Data:         nil,
	}, binfo)
	assert.Nil(t, binfo.Result)
	// Ensure AfterRemoteHooks was called
	ainfo := h.afterRemoteCallInfo
	assertRemoteInfo(t, &RemoteHooksInfo{
		Endpoint:     "List",
		ResourceName: "Person",
		ResourceID:   nil,
		Data:         nil,
		Result:       nil,
	}, ainfo)
	assert.NotNil(t, ainfo.Result)
	var hres Res
	utils.DecodeMapToStruct(ainfo.Result, &hres)
	r.Equal(res.Count, hres.Count)
}

func TestBeforeRemoteHooksListError(t *testing.T) {
	r := require.New(t)
	k := NewTestKit(t)
	h := RemoteHooksMockError{}
	k.linker.AddResource(Resource{
		Name:        "Person",
		Datasource:  gormds.NewGormDatasource(k.db, &testutils.Person{}),
		RemoteHooks: &h,
	})
	w := httptest.NewRecorder()
	k.app.InitApp()
	k.app.ServeHTTP(w, httptest.NewRequest("GET", "/api/people", nil))
	r.Equal(http.StatusBadRequest, w.Code)
}

func TestBeforeRemoteHooksShow(t *testing.T) {
	r := require.New(t)
	k := NewTestKit(t)
	h := RemoteHooksMock{}
	k.linker.AddResource(Resource{
		Name:        "Person",
		Datasource:  gormds.NewGormDatasource(k.db, &testutils.Person{}),
		RemoteHooks: &h,
	})
	w := httptest.NewRecorder()
	k.app.InitApp()
	k.app.ServeHTTP(w, httptest.NewRequest("GET", "/api/people/1", nil))
	type Res struct {
		testutils.Person
		CustomField string `json:"custom_field"`
	}
	var res Res
	json.Unmarshal(w.Body.Bytes(), &res)
	r.Equal(http.StatusOK, w.Code)
	r.Equal(uint(1), res.ID)
	r.Equal("Carlos Strand", res.Name)
	r.Equal("perfect!", res.CustomField)
	// Ensure BeforeRemoteHooks was called
	binfo := h.beforeRemoteCallInfo
	assertRemoteInfo(t, &RemoteHooksInfo{
		Endpoint:     "Show",
		ResourceName: "Person",
		ResourceID:   ptr.String("1"),
		Data:         nil,
	}, binfo)
	assert.Nil(t, binfo.Result)
	// Ensure AfterRemoteHooks was called
	ainfo := h.afterRemoteCallInfo
	assertRemoteInfo(t, &RemoteHooksInfo{
		Endpoint:     "Show",
		ResourceName: "Person",
		ResourceID:   ptr.String("1"),
		Data:         nil,
	}, ainfo)
	assert.NotNil(t, ainfo.Result)
	var hres Res
	utils.DecodeMapToStruct(ainfo.Result, &hres)
	r.Equal(res.ID, hres.ID)
	r.Equal(res.Name, hres.Name)
}

func TestBeforeRemoteHooksShowError(t *testing.T) {
	r := require.New(t)
	k := NewTestKit(t)
	h := RemoteHooksMockError{}
	k.linker.AddResource(Resource{
		Name:        "Person",
		Datasource:  gormds.NewGormDatasource(k.db, &testutils.Person{}),
		RemoteHooks: &h,
	})
	w := httptest.NewRecorder()
	k.app.InitApp()
	k.app.ServeHTTP(w, httptest.NewRequest("GET", "/api/people/1", nil))
	r.Equal(http.StatusBadRequest, w.Code)
}

func TestBeforeRemoteHooksCreate(t *testing.T) {
	r := require.New(t)
	k := NewTestKit(t)
	h := RemoteHooksMock{}
	k.linker.AddResource(Resource{
		Name:        "Person",
		Datasource:  gormds.NewGormDatasource(k.db, &testutils.Person{}),
		RemoteHooks: &h,
	})
	w := httptest.NewRecorder()
	k.app.InitApp()
	body := `
		{
			"name": "Bruce Wayne",
			"email": "bruce@test.com"
		}
	`
	req := httptest.NewRequest("POST", "/api/people", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	k.app.ServeHTTP(w, req)
	type Res struct {
		testutils.Person
		CustomField string `json:"custom_field"`
	}
	var res Res
	json.Unmarshal(w.Body.Bytes(), &res)
	r.Equal(http.StatusOK, w.Code)
	r.Equal(uint(4), res.ID)
	r.Equal("Bruce Wayne", res.Name)
	r.Equal("perfect!", res.CustomField)
	// Ensure BeforeRemoteHooks was called
	binfo := h.beforeRemoteCallInfo
	assertRemoteInfo(t, &RemoteHooksInfo{
		Endpoint:     "Create",
		ResourceName: "Person",
		Data:         &map[string]interface{}{"email": "bruce@test.com", "name": "Bruce Wayne"},
	}, binfo)
	assert.Nil(t, binfo.Result)
	// Ensure AfterRemoteHooks was called
	ainfo := h.afterRemoteCallInfo
	assertRemoteInfo(t, &RemoteHooksInfo{
		Endpoint:     "Create",
		ResourceName: "Person",
		ResourceID:   ptr.String("4"),
		Data:         &map[string]interface{}{"email": "bruce@test.com", "name": "Bruce Wayne"},
	}, ainfo)
	assert.NotNil(t, ainfo.Result)
	var hres Res
	utils.DecodeMapToStruct(ainfo.Result, &hres)
	r.Equal(res.ID, hres.ID)
	r.Equal(res.Name, hres.Name)
}

func TestBeforeRemoteHooksCreateError(t *testing.T) {
	r := require.New(t)
	k := NewTestKit(t)
	h := RemoteHooksMockError{}
	k.linker.AddResource(Resource{
		Name:        "Person",
		Datasource:  gormds.NewGormDatasource(k.db, &testutils.Person{}),
		RemoteHooks: &h,
	})
	w := httptest.NewRecorder()
	k.app.InitApp()
	body := `
		{
			"name": "Bruce Wayne",
			"email": "bruce@test.com"
		}
	`
	k.app.ServeHTTP(w, httptest.NewRequest("POST", "/api/people", strings.NewReader(body)))
	r.Equal(http.StatusBadRequest, w.Code)
}

func TestBeforeRemoteHooksUpdate(t *testing.T) {
	r := require.New(t)
	k := NewTestKit(t)
	h := RemoteHooksMock{}
	k.linker.AddResource(Resource{
		Name:        "Person",
		Datasource:  gormds.NewGormDatasource(k.db, &testutils.Person{}),
		RemoteHooks: &h,
	})
	w := httptest.NewRecorder()
	k.app.InitApp()
	body := `
		{
			"name": "Bruce Wayne",
			"email": "bruce@test.com"
		}
	`
	req := httptest.NewRequest("PUT", "/api/people/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	k.app.ServeHTTP(w, req)
	type Res struct {
		testutils.Person
		CustomField string `json:"custom_field"`
	}
	var res Res
	json.Unmarshal(w.Body.Bytes(), &res)
	r.Equal(http.StatusOK, w.Code)
	r.Equal(uint(1), res.ID)
	r.Equal("Bruce Wayne", res.Name)
	r.Equal("perfect!", res.CustomField)
	// Ensure BeforeRemoteHooks was called
	binfo := h.beforeRemoteCallInfo
	assertRemoteInfo(t, &RemoteHooksInfo{
		Endpoint:     "Update",
		ResourceName: "Person",
		ResourceID:   ptr.String("1"),
		Data:         &map[string]interface{}{"email": "bruce@test.com", "name": "Bruce Wayne"},
	}, binfo)
	assert.Nil(t, binfo.Result)
	// Ensure AfterRemoteHooks was called
	ainfo := h.afterRemoteCallInfo
	assertRemoteInfo(t, &RemoteHooksInfo{
		Endpoint:     "Update",
		ResourceName: "Person",
		ResourceID:   ptr.String("1"),
		Data:         &map[string]interface{}{"email": "bruce@test.com", "name": "Bruce Wayne"},
	}, ainfo)
	assert.NotNil(t, ainfo.Result)
	var hres Res
	utils.DecodeMapToStruct(ainfo.Result, &hres)
	r.Equal(res.ID, hres.ID)
	r.Equal(res.Name, hres.Name)
}

func TestBeforeRemoteHooksUpdateError(t *testing.T) {
	r := require.New(t)
	k := NewTestKit(t)
	h := RemoteHooksMockError{}
	k.linker.AddResource(Resource{
		Name:        "Person",
		Datasource:  gormds.NewGormDatasource(k.db, &testutils.Person{}),
		RemoteHooks: &h,
	})
	w := httptest.NewRecorder()
	k.app.InitApp()
	body := `
		{
			"name": "Bruce Wayne",
			"email": "bruce@test.com"
		}
	`
	k.app.ServeHTTP(w, httptest.NewRequest("PUT", "/api/people/1", strings.NewReader(body)))
	r.Equal(http.StatusBadRequest, w.Code)
}

func TestBeforeRemoteHooksDestroy(t *testing.T) {
	r := require.New(t)
	k := NewTestKit(t)
	h := RemoteHooksMock{}
	k.linker.AddResource(Resource{
		Name:        "Person",
		Datasource:  gormds.NewGormDatasource(k.db, &testutils.Person{}),
		RemoteHooks: &h,
	})
	w := httptest.NewRecorder()
	k.app.InitApp()
	k.app.ServeHTTP(w, httptest.NewRequest("DELETE", "/api/people/1", nil))
	type DelRes struct {
		Deleted     bool   `json:"deleted"`
		CustomField string `json:"custom_field"`
	}
	var res DelRes
	json.Unmarshal(w.Body.Bytes(), &res)
	r.Equal(http.StatusOK, w.Code)
	r.Equal(true, res.Deleted)
	r.Equal("perfect!", res.CustomField)
	// Ensure BeforeRemoteHooks was called
	binfo := h.beforeRemoteCallInfo
	assertRemoteInfo(t, &RemoteHooksInfo{
		Endpoint:     "Destroy",
		ResourceName: "Person",
		ResourceID:   ptr.String("1"),
		Data:         nil,
	}, binfo)
	assert.Nil(t, binfo.Result)
	// Ensure AfterRemoteHooks was called
	ainfo := h.afterRemoteCallInfo
	assertRemoteInfo(t, &RemoteHooksInfo{
		Endpoint:     "Destroy",
		ResourceName: "Person",
		ResourceID:   ptr.String("1"),
		Data:         nil,
	}, ainfo)
	assert.NotNil(t, ainfo.Result)
	var hres DelRes
	utils.DecodeMapToStruct(ainfo.Result, &hres)
	r.Equal(res.Deleted, hres.Deleted)
}
