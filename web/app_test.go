package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-zepto/zepto/testutils"
)

func CreateAppTest() *App {
	app := NewApp()
	return app
}

func TestApp_NewApp(t *testing.T) {
	app := NewApp()
	if app.muxRouter == nil {
		t.Errorf("app.muxRouter should not be nil")
	}
	if app.opts.tmplEngine == nil {
		t.Errorf("app.opts.tmplEngine should not be nil")
	}
}

func setupAppTest() *App {
	app := CreateAppTest()
	app.tmplEngine = &testutils.EngineMock{}
	app.opts.env = "development"
	return app
}

func assertRequestStatus(t *testing.T, app *App, method string, path string, code int) {
	w := httptest.NewRecorder()
	app.Init()
	app.ServeHTTP(w, httptest.NewRequest(method, path, nil))
	if w.Code != code {
		t.Error("Did not get expected HTTP status code, got", w.Code)
	}
}

func assertRequest(t *testing.T, app *App, method string, path string, code int, body string) {
	w := httptest.NewRecorder()
	app.Init()
	app.ServeHTTP(w, httptest.NewRequest(method, path, nil))
	if w.Code != code {
		t.Error("Did not get expected HTTP status code, got", w.Code)
	}
	if w.Body.String() != body {
		t.Error("Did not get expected body, got", w.Body.String())
	}
}

func assertRequestFn(t *testing.T, app *App, method string, path string, fn func(w *httptest.ResponseRecorder)) {
	w := httptest.NewRecorder()
	app.ServeHTTP(w, httptest.NewRequest(method, path, nil))
	fn(w)
}

var DefaultRouterHandler = func(ctx Context) error {
	ctx.Render("hello/index")
	return nil
}

func TestApp_HandleMethod_Get(t *testing.T) {
	app := setupAppTest()
	app.Get("/hello", DefaultRouterHandler)
	app.Init()
	assertRequest(t, app, "GET", "/hello", 200, "Mocked Template!")
}

func TestApp_HandleMethod_Post(t *testing.T) {
	app := setupAppTest()
	app.Post("/hello", DefaultRouterHandler)
	app.Init()
	assertRequest(t, app, "POST", "/hello", 200, "Mocked Template!")
}

func TestApp_HandleMethod_Put(t *testing.T) {
	app := setupAppTest()
	app.Put("/hello", DefaultRouterHandler)
	app.Init()
	assertRequest(t, app, "PUT", "/hello", 200, "Mocked Template!")
}

func TestApp_HandleMethod_Delete(t *testing.T) {
	app := setupAppTest()
	app.Delete("/hello", DefaultRouterHandler)
	assertRequest(t, app, "DELETE", "/hello", 200, "Mocked Template!")
}

func TestApp_HandleMethod_Patch(t *testing.T) {
	app := setupAppTest()
	app.Patch("/hello", DefaultRouterHandler)
	assertRequest(t, app, "PATCH", "/hello", 200, "Mocked Template!")
}

func TestApp_HandleMethod_ControllerPanicDevelopmentAsDebug(t *testing.T) {
	app := setupAppTest()
	app.Get("/hello", func(ctx Context) error {
		panic(errors.New("panic problem"))
		return nil
	})
	w := httptest.NewRecorder()
	app.Init()
	app.ServeHTTP(w, httptest.NewRequest("GET", "/hello", nil))
	if w.Code != http.StatusInternalServerError {
		t.Error("Did not get expected HTTP status code, got", w.Code)
	}
	if !strings.Contains(w.Body.String(), "<title>Internal Error - 500 - internal server error: panic problem</title>") {
		t.Error("Did not get expected body, got", w.Body.String())
	}
}

func assertProdError(t *testing.T, handler RouteHandler) {
	app := setupAppTest()
	app.opts.env = "production"
	app.Get("/hello", handler)
	w := httptest.NewRecorder()
	app.Init()
	app.ServeHTTP(w, httptest.NewRequest("GET", "/hello", nil))
	if w.Code != http.StatusInternalServerError {
		t.Error("Did not get expected HTTP status code, got", w.Code)
	}
	if w.Body.String() != "500 - internal server error" {
		t.Error("Did not get expected body, got", w.Body.String())
	}
}

func TestApp_HandleMethod_ControllerErrorProduction(t *testing.T) {
	os.Setenv("SESSION_SECRET", "test_session_secret")
	assertProdError(t, func(ctx Context) error {
		return errors.New("some error in prod")
	})
	assertProdError(t, func(ctx Context) error {
		panic(errors.New("some panic error in prod"))
		return nil
	})
	assertProdError(t, func(ctx Context) error {
		panic("some panic in prod")
		return nil
	})
}

func TestApp_Start(t *testing.T) {
	tmplEngine := &testutils.EngineMock{}
	app := setupAppTest()
	app.opts.env = "production"
	app.opts.webpackEnabled = false
	app.tmplEngine = tmplEngine
	if tmplEngine.InitCalled {
		t.Errorf("Init should not be called at this point")
	}
	app.Init()
	if !tmplEngine.InitCalled {
		t.Errorf("Init should be called at this point")
	}
}

func TestApp_HandleMethod_ControllerRenderJson(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	app := setupAppTest()
	p := &Person{"Mike", 27}
	app.Get("/person", func(ctx Context) error {
		return ctx.RenderJson(p)
	})
	w := httptest.NewRecorder()
	app.Init()
	app.ServeHTTP(w, httptest.NewRequest("GET", "/person", nil))
	if w.Code != http.StatusOK {
		t.Error("Did not get expected HTTP status code, got", w.Code)
	}
	var pp Person
	err := json.Unmarshal(w.Body.Bytes(), &pp)
	if err != nil {
		t.Error(err)
	}
	if pp.Name != "Mike" || pp.Age != 27 {
		t.Errorf("Decoded json data is not as expected. Received: %#v", pp)
	}
}
