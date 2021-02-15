package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func SetupTestRouterApi(t *testing.T) *Router {
	r := require.New(t)
	app := NewApp()
	apiv1Router := app.Router("/api/v1")
	r.NotNil(apiv1Router)
	r.Len(app.routers, 1)
	r.Equal(app.routers[0].options.path, "/api/v1")
	return apiv1Router
}

func TestNewRouter(t *testing.T) {
	SetupTestRouterApi(t)
}

func TestNewRouterWithRoutes(t *testing.T) {
	r := require.New(t)
	apiv1Router := SetupTestRouterApi(t)
	r.Len(apiv1Router.handlers, 0)
	apiv1Router.HandleMethod([]string{"GET"}, "/hello", func(ctx Context) error {
		return ctx.RenderJson(map[string]string{"hello": "world"})
	})
	r.Len(apiv1Router.handlers, 1)
	h := apiv1Router.handlers[0]
	r.Equal(h.methods, []string{"GET"})
	r.Equal(h.path, "/hello")
}

func TestRouter_ManyRouters(t *testing.T) {
	r := require.New(t)
	type Info struct {
		Version string `json:"version"`
	}
	app := setupAppTest()

	apiv1Router := app.Router("/api/v1")
	apiv1Router.Get("/info", func(ctx Context) error {
		return ctx.RenderJson(Info{Version: "v1"})
	})

	apiv2Router := app.Router("/api/v2")
	apiv2Router.Get("/info", func(ctx Context) error {
		return ctx.RenderJson(Info{Version: "v2"})
	})

	apiv3Router := app.Router("/api/v3")
	apiv3Router.Get("/info", func(ctx Context) error {
		return ctx.RenderJson(Info{Version: "v3"})
	})

	app.Start()

	for _, v := range []string{"v1", "v2", "v3"} {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, httptest.NewRequest("GET", "/api/"+v+"/info", nil))
		if w.Code != http.StatusOK {
			t.Error("Did not get expected HTTP status code, got", w.Code)
		}
		var info Info
		json.Unmarshal(w.Body.Bytes(), &info)
		r.Equal(v, info.Version)
	}
}

func TestRouter_ManyRoutersWithHosts(t *testing.T) {
	r := require.New(t)
	type Info struct {
		Host string `json:"host"`
	}
	app := setupAppTest()

	router1 := app.Router("/api", Hosts("go-zepto.com"))
	router1.Get("/info", func(ctx Context) error {
		return ctx.RenderJson(Info{Host: "go-zepto.com"})
	})

	router2 := app.Router("/api", Hosts("go-zepto.ca"))
	router2.Get("/info", func(ctx Context) error {
		return ctx.RenderJson(Info{Host: "go-zepto.ca"})
	})

	app.Start()

	for _, host := range []string{"go-zepto.com", "go-zepto.ca"} {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/info", nil)
		req.Host = host
		app.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Error("Did not get expected HTTP status code, got", w.Code)
		}
		var info Info
		json.Unmarshal(w.Body.Bytes(), &info)
		r.Equal(host, info.Host)
	}
}

func TestRouter_MultipleHosts(t *testing.T) {
	type Info struct {
		Message string `json:"message"`
	}
	app := setupAppTest()

	router1 := app.Router("/api", Hosts("go-zepto.com", "go-zepto.ca"))
	router1.Get("/info", func(ctx Context) error {
		return ctx.RenderJson(Info{Message: "Hello World"})
	})

	app.Start()

	var cases = []struct {
		Host           string
		ExpectedStatus int
	}{
		{
			Host:           "go-zepto.com",
			ExpectedStatus: 200,
		},
		{
			Host:           "go-zepto.ca",
			ExpectedStatus: 200,
		},
		{
			Host:           "go-zepto.it",
			ExpectedStatus: 404,
		},
	}
	for _, c := range cases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/info", nil)
		req.Host = c.Host
		app.ServeHTTP(w, req)
		if w.Code != c.ExpectedStatus {
			t.Error("Did not get expected HTTP status code, got", w.Code)
		}
	}
}

type ResourceMockInfo struct {
	Controller string `json:"controller"`
}

type ResourceMock struct{}

func (rm *ResourceMock) List(ctx Context) error {
	return ctx.RenderJson(ResourceMockInfo{Controller: "LIST BOOKS"})
}

func (rm *ResourceMock) Show(ctx Context) error {
	id := ctx.Params()["id"]
	return ctx.RenderJson(ResourceMockInfo{Controller: "SHOW BOOK ID=" + id})
}

func (rm *ResourceMock) Create(ctx Context) error {
	return ctx.RenderJson(ResourceMockInfo{Controller: "CREATE BOOK"})
}

func (rm *ResourceMock) Update(ctx Context) error {
	id := ctx.Params()["id"]
	return ctx.RenderJson(ResourceMockInfo{Controller: "UPDATE BOOK ID=" + id})
}

func (rm *ResourceMock) Destroy(ctx Context) error {
	id := ctx.Params()["id"]
	return ctx.RenderJson(ResourceMockInfo{Controller: "DESTROY BOOK ID=" + id})
}

func TestRouter_Resource(t *testing.T) {
	r := require.New(t)
	app := setupAppTest()
	api := app.Router("/api")
	api.Resource("/books", &ResourceMock{})
	app.Start()
	testCases := []struct {
		method       string
		endpoint     string
		expectedText string
	}{
		{
			method:       "GET",
			endpoint:     "/api/books",
			expectedText: "LIST BOOKS",
		},
		{
			method:       "GET",
			endpoint:     "/api/books/55",
			expectedText: "SHOW BOOK ID=55",
		},
		{
			method:       "POST",
			endpoint:     "/api/books",
			expectedText: "CREATE BOOK",
		},
		{
			method:       "PUT",
			endpoint:     "/api/books/90",
			expectedText: "UPDATE BOOK ID=90",
		},
		{
			method:       "DELETE",
			endpoint:     "/api/books/70",
			expectedText: "DESTROY BOOK ID=70",
		},
	}
	for _, c := range testCases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(c.method, c.endpoint, nil)
		app.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Error("Did not get expected HTTP status code, got", w.Code)
		}
		var info ResourceMockInfo
		json.Unmarshal(w.Body.Bytes(), &info)
		r.Equal(c.expectedText, info.Controller)
	}
}
