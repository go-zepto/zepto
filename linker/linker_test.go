package linker

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gormds "github.com/go-zepto/zepto/linker/datasource/gorm"
	"github.com/go-zepto/zepto/linker/datasource/gorm/testutils"
	"github.com/go-zepto/zepto/linker/hooks"
	"github.com/go-zepto/zepto/linker/utils"
	"github.com/go-zepto/zepto/web"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type TestKit struct {
	app    *web.App
	router *web.Router
	db     *gorm.DB
	linker *Linker
}

func NewTestKit(t *testing.T) TestKit {
	r := require.New(t)
	app := web.NewApp()
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
	beforeRemoteCallInfo *hooks.RemoteHooksInfo
	afterRemoteCallInfo  *hooks.RemoteHooksInfo
}

func (r *RemoteHooksMock) BeforeRemote(info hooks.RemoteHooksInfo) error {
	r.beforeRemoteCallInfo = &info
	return nil
}

func (r *RemoteHooksMock) AfterRemote(info hooks.RemoteHooksInfo) error {
	r.afterRemoteCallInfo = &info
	return nil
}

type RemoteHooksMockError struct{}

func (r *RemoteHooksMockError) BeforeRemote(info hooks.RemoteHooksInfo) error {
	info.Ctx.SetStatus(400)
	return errors.New("ups! strange error")
}

func (r *RemoteHooksMockError) AfterRemote(info hooks.RemoteHooksInfo) error {
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
	k.app.Start()
	k.app.ServeHTTP(w, httptest.NewRequest("GET", "/api/people", nil))
	type Res struct {
		Data  []testutils.Person `json:"data"`
		Count int                `json:"count"`
	}
	var res Res
	json.Unmarshal(w.Body.Bytes(), &res)
	r.Equal(http.StatusOK, w.Code)
	r.Equal(3, res.Count)
	r.Equal("Carlos Strand", res.Data[0].Name)
	r.Equal("Bill Gates", res.Data[1].Name)
	r.Equal("Clark Kent", res.Data[2].Name)
	// Ensure BeforeRemoteHooks was called
	binfo := h.beforeRemoteCallInfo
	r.NotNil(binfo)
	r.NotNil(binfo.Ctx)
	r.Equal("List", binfo.Endpoint)
	r.Nil(binfo.Data)
	// Ensure AfterRemoteHooks was called
	ainfo := h.afterRemoteCallInfo
	r.NotNil(ainfo)
	r.NotNil(ainfo.Data)
	r.Equal("List", ainfo.Endpoint)
	var hres Res
	utils.DecodeMapToStruct(ainfo.Data, &hres)
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
	k.app.Start()
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
	k.app.Start()
	k.app.ServeHTTP(w, httptest.NewRequest("GET", "/api/people/1", nil))
	var res testutils.Person
	json.Unmarshal(w.Body.Bytes(), &res)
	r.Equal(http.StatusOK, w.Code)
	r.Equal(uint(1), res.ID)
	r.Equal("Carlos Strand", res.Name)
	// Ensure BeforeRemoteHooks was called
	binfo := h.beforeRemoteCallInfo
	r.NotNil(binfo)
	r.NotNil(binfo.Ctx)
	r.Equal("Show", binfo.Endpoint)
	r.Nil(binfo.Data)
	// Ensure AfterRemoteHooks was called
	ainfo := h.afterRemoteCallInfo
	r.NotNil(ainfo)
	r.NotNil(ainfo.Data)
	r.Equal("Show", ainfo.Endpoint)
	r.NotNil(ainfo.ID)
	r.Equal("1", *ainfo.ID)
	var hres testutils.Person
	utils.DecodeMapToStruct(ainfo.Data, &hres)
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
	k.app.Start()
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
	k.app.Start()
	body := `
		{
			"name": "Bruce Wayne",
			"email": "bruce@test.com"
		}
	`
	k.app.ServeHTTP(w, httptest.NewRequest("POST", "/api/people", strings.NewReader(body)))
	var res testutils.Person
	json.Unmarshal(w.Body.Bytes(), &res)
	r.Equal(http.StatusOK, w.Code)
	r.Equal(uint(4), res.ID)
	r.Equal("Bruce Wayne", res.Name)
	// Ensure BeforeRemoteHooks was called
	binfo := h.beforeRemoteCallInfo
	r.NotNil(binfo)
	r.NotNil(binfo.Ctx)
	r.Equal("Create", binfo.Endpoint)
	r.NotNil(binfo.Data)
	binfoData := *binfo.Data
	r.Equal("Bruce Wayne", binfoData["name"])
	r.Equal("bruce@test.com", binfoData["email"])
	// Ensure AfterRemoteHooks was called
	ainfo := h.afterRemoteCallInfo
	r.NotNil(ainfo)
	r.NotNil(ainfo.Data)
	r.NotNil(ainfo.ID)
	r.Equal("4", *ainfo.ID)
	r.Equal("Create", ainfo.Endpoint)
	var hres testutils.Person
	utils.DecodeMapToStruct(ainfo.Data, &hres)
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
	k.app.Start()
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
	k.app.Start()
	body := `
		{
			"name": "Bruce Wayne",
			"email": "bruce@test.com"
		}
	`
	k.app.ServeHTTP(w, httptest.NewRequest("PUT", "/api/people/1", strings.NewReader(body)))
	var res testutils.Person
	json.Unmarshal(w.Body.Bytes(), &res)
	r.Equal(http.StatusOK, w.Code)
	r.Equal(uint(1), res.ID)
	r.Equal("Bruce Wayne", res.Name)
	// Ensure BeforeRemoteHooks was called
	binfo := h.beforeRemoteCallInfo
	r.NotNil(binfo)
	r.NotNil(binfo.Ctx)
	r.Equal("Update", binfo.Endpoint)
	r.NotNil(binfo.Data)
	binfoData := *binfo.Data
	r.Equal("Bruce Wayne", binfoData["name"])
	r.Equal("bruce@test.com", binfoData["email"])
	// Ensure AfterRemoteHooks was called
	ainfo := h.afterRemoteCallInfo
	r.NotNil(ainfo)
	r.NotNil(ainfo.Data)
	r.NotNil(ainfo.ID)
	r.Equal("1", *ainfo.ID)
	r.Equal("Update", ainfo.Endpoint)
	var hres testutils.Person
	utils.DecodeMapToStruct(ainfo.Data, &hres)
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
	k.app.Start()
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
	k.app.Start()
	k.app.ServeHTTP(w, httptest.NewRequest("DELETE", "/api/people/1", nil))
	type DelRes struct {
		Deleted bool `json:"deleted"`
	}
	var res DelRes
	json.Unmarshal(w.Body.Bytes(), &res)
	r.Equal(http.StatusOK, w.Code)
	r.Equal(true, res.Deleted)
	// Ensure BeforeRemoteHooks was called
	binfo := h.beforeRemoteCallInfo
	r.NotNil(binfo)
	r.NotNil(binfo.Ctx)
	r.Equal("Destroy", binfo.Endpoint)
	r.Nil(binfo.Data)
	// Ensure AfterRemoteHooks was called
	ainfo := h.afterRemoteCallInfo
	r.NotNil(ainfo)
	r.NotNil(ainfo.Data)
	r.Equal("Destroy", ainfo.Endpoint)
	r.NotNil(ainfo.ID)
	r.Equal("1", *ainfo.ID)
	var hres DelRes
	utils.DecodeMapToStruct(ainfo.Data, &hres)
	r.Equal(res.Deleted, hres.Deleted)
}
