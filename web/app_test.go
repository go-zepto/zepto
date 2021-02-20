package web

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-zepto/zepto/broker"
	"github.com/go-zepto/zepto/logger/logrus"
	"github.com/go-zepto/zepto/testutils"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func CreateAppWithBrokerTest() (*App, *testutils.BrokerProviderMock) {
	l := logrus.NewLogrus(log.New())
	bp := &testutils.BrokerProviderMock{}
	b := broker.NewBroker(l, bp)
	app := NewApp(
		Broker(b),
	)
	return app, bp
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

func TestApp_Broker_Init(t *testing.T) {
	app, bp := CreateAppWithBrokerTest()
	if app.opts.broker == nil {
		t.Fatal("app broker should not be nil")
	}
	if bp.InitCalled {
		t.Fatal("broker init should not be called at this moment")
	}
	bp.Init(&broker.InitOptions{
		Logger: app.opts.logger,
	})
	if !bp.InitCalled {
		t.Fatal("broker init should be called at this moment")
	}
}

func TestApp_Broker_Publish(t *testing.T) {
	_, bp := CreateAppWithBrokerTest()
	if bp.PublishCalled {
		t.Fatal("broker publish should not be called at this moment")
	}
	msg := &broker.Message{
		Header: map[string]string{"hello": "world"},
		Body:   []byte("[]"),
	}
	err := bp.Publish(context.Background(), "my.topic", msg)
	if err != nil {
		t.Fatal(err)
	}
	if !bp.PublishCalled {
		t.Fatal("broker publish should be called at this moment")
	}
	if bp.PublishCallArgs.Topic != "my.topic" {
		t.Fatal("broker publish should have topic=my.topic")
	}
	if bp.PublishCallArgs.Msg != msg {
		t.Errorf("broker publish should have message=%s", msg)
	}
}

func TestApp_Broker_Subscribe(t *testing.T) {
	_, bp := CreateAppWithBrokerTest()
	if bp.SubscribeCalled {
		t.Fatal("broker subscribe should not be called at this moment")
	}
	handler := func(ctx context.Context, msg *broker.Message) {}
	err := bp.Subscribe(context.Background(), "my.topic", handler)
	if err != nil {
		t.Fatal(err)
	}
	if !bp.SubscribeCalled {
		t.Fatal("broker subscribe should be called at this moment")
	}
	if bp.SubscribeCallArgs.Topic != "my.topic" {
		t.Fatal("broker subscribe should have topic=my.topic")
	}
}

func setupAppTest() *App {
	app, _ := CreateAppWithBrokerTest()
	app.opts.env = "development"
	muxRouter := mux.NewRouter()
	app.muxRouter = muxRouter
	app.tmplEngine = &testutils.EngineMock{}
	app.Init(InitOptions{})
	return app
}

func assertRequestStatus(t *testing.T, app *App, method string, path string, code int) {
	w := httptest.NewRecorder()
	app.ServeHTTP(w, httptest.NewRequest(method, path, nil))
	if w.Code != code {
		t.Error("Did not get expected HTTP status code, got", w.Code)
	}
}

func assertRequest(t *testing.T, app *App, method string, path string, code int, body string) {
	w := httptest.NewRecorder()
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
	app.Init(InitOptions{})
	assertRequest(t, app, "GET", "/hello", 200, "Mocked Template!")
}

func TestApp_HandleMethod_Post(t *testing.T) {
	app := setupAppTest()
	app.Post("/hello", DefaultRouterHandler)
	app.Init(InitOptions{})
	assertRequest(t, app, "POST", "/hello", 200, "Mocked Template!")
}

func TestApp_HandleMethod_Put(t *testing.T) {
	app := setupAppTest()
	app.Put("/hello", DefaultRouterHandler)
	app.Init(InitOptions{})
	assertRequest(t, app, "PUT", "/hello", 200, "Mocked Template!")
}

func TestApp_HandleMethod_Delete(t *testing.T) {
	app := setupAppTest()
	app.Delete("/hello", DefaultRouterHandler)
	app.Init(InitOptions{})
	assertRequest(t, app, "DELETE", "/hello", 200, "Mocked Template!")
}

func TestApp_HandleMethod_Patch(t *testing.T) {
	app := setupAppTest()
	app.Patch("/hello", DefaultRouterHandler)
	app.Init(InitOptions{})
	assertRequest(t, app, "PATCH", "/hello", 200, "Mocked Template!")
}

func TestApp_HandleMethod_ControllerPanicDevelopmentAsDebug(t *testing.T) {
	app := setupAppTest()
	app.Get("/hello", func(ctx Context) error {
		panic(errors.New("panic problem"))
		return nil
	})
	app.Init(InitOptions{})
	w := httptest.NewRecorder()
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
	app.Init(InitOptions{})
	w := httptest.NewRecorder()
	app.ServeHTTP(w, httptest.NewRequest("GET", "/hello", nil))
	if w.Code != http.StatusInternalServerError {
		t.Error("Did not get expected HTTP status code, got", w.Code)
	}
	if w.Body.String() != "500 - internal server error" {
		t.Error("Did not get expected body, got", w.Body.String())
	}
}

func TestApp_HandleMethod_ControllerErrorProduction(t *testing.T) {
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
	app.Init(InitOptions{})
	app.Start()
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
	app.Init(InitOptions{})
	w := httptest.NewRecorder()
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
