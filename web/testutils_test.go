package web

import (
	"fmt"
	"testing"
)

func HelloHandler(ctx Context) error {
	ctx.SetStatus(201)
	return ctx.RenderJson(map[string]string{"hello": "world"})
}

func TmplHandler(ctx Context) error {
	ctx.SetStatus(201)
	return ctx.Render("hello")
}

func TestHandlerRequestJson(t *testing.T) {
	app := NewApp()
	zt := NewZeptoTest(t, app)
	res, err := zt.TestHandlerRequest(HelloHandler, "GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	res.AssertStatusCode(201)
	res.AssertBodyEquals("{\"hello\":\"world\"}\n")
	fmt.Println(res)
}

func TestHandlerRequestTemplate(t *testing.T) {
	app := NewApp()
	zt := NewZeptoTest(t, app)
	res, err := zt.TestHandlerRequest(TmplHandler, "GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	res.AssertStatusCode(201)
	res.AssertBodyContains("This is a template used for tests")
}
