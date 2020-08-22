package web

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHandlerRequestJson(t *testing.T) {
	var HelloHandler = func(ctx Context) error {
		return ctx.RenderJson(map[string]string{"hello": "world"})
	}
	app := NewApp()
	zt := NewZeptoTest(t, app)
	res, err := zt.TestHandlerRequest(TestHandlerRequestOptions{
		Handler: HelloHandler,
		Method:  "GET",
		Target:  "/",
		Body:    nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	res.AssertBodyEquals("{\"hello\":\"world\"}\n")
}

func TestHandlerRequestTemplate(t *testing.T) {
	var TmplHandler = func(ctx Context) error {
		return ctx.Render("hello")
	}
	app := NewApp()
	zt := NewZeptoTest(t, app)
	res, err := zt.TestHandlerRequest(TestHandlerRequestOptions{
		Handler: TmplHandler,
		Method:  "GET",
		Target:  "/",
		InitialSession: map[string]string{
			"abc": "123",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	res.AssertBodyContains("This is a template used for tests")
}

func TestHandlerRequestStatus(t *testing.T) {
	var handler = func(ctx Context) error {
		ctx.SetStatus(401)
		return ctx.RenderJson(map[string]string{"hello": "world"})
	}
	app := NewApp()
	zt := NewZeptoTest(t, app)
	res, err := zt.TestHandlerRequest(TestHandlerRequestOptions{
		Handler: handler,
		Method:  "GET",
		Target:  "/",
	})
	if err != nil {
		t.Fatal(err)
	}
	res.AssertStatusCode(401)
}

func TestHandlerRequestHeaderReq(t *testing.T) {
	var handler = func(ctx Context) error {
		assert.Equal(t, "OK", ctx.Request().Header.Get("X-TESUTILS-HEADER"))
		return ctx.RenderJson(map[string]string{"hello": "world"})
	}
	app := NewApp()
	zt := NewZeptoTest(t, app)
	header := http.Header{}
	header.Set("X-TESUTILS-HEADER", "OK")
	_, err := zt.TestHandlerRequest(TestHandlerRequestOptions{
		Handler: handler,
		Method:  "GET",
		Target:  "/",
		Header:  header,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandlerRequestHeaderRes(t *testing.T) {
	var handler = func(ctx Context) error {
		ctx.Response().Header().Set("X-HELLO", "WORLD")
		return ctx.RenderJson(map[string]string{"hello": "world"})
	}
	app := NewApp()
	zt := NewZeptoTest(t, app)
	res, err := zt.TestHandlerRequest(TestHandlerRequestOptions{
		Handler: handler,
		Method:  "GET",
		Target:  "/",
	})
	if err != nil {
		t.Fatal(err)
	}
	res.AssertHeaderValue("X-HELLO", "WORLD")
}

func TestHandlerRequestHost(t *testing.T) {
	var handler = func(ctx Context) error {
		assert.Equal(t, "go-zepto.github.io", ctx.Request().Host)
		return ctx.RenderJson(map[string]string{"hello": "world"})
	}
	app := NewApp()
	zt := NewZeptoTest(t, app)
	_, err := zt.TestHandlerRequest(TestHandlerRequestOptions{
		Handler: handler,
		Method:  "GET",
		Target:  "/",
		Host:    "go-zepto.github.io",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandlerRequestSession(t *testing.T) {
	var handler = func(ctx Context) error {
		ctx.Session().Set("user_id", "42")
		return ctx.RenderJson(map[string]string{"hello": "world"})
	}
	app := NewApp()
	zt := NewZeptoTest(t, app)
	res, err := zt.TestHandlerRequest(TestHandlerRequestOptions{
		Handler: handler,
		Method:  "GET",
		Target:  "/",
		InitialSession: map[string]string{
			"initial_session": "ok",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	res.AssertSessionValue("initial_session", "ok")
	res.AssertSessionValue("user_id", "42")
}
