package web

import (
	"testing"
)

type MiddlewareTester struct{}

func TestMiddlewareStack_Use_Default(t *testing.T) {
	var Middleware1 = func(next RouteHandler) RouteHandler {
		return func(ctx Context) error {
			ctx.Set("middleware-order", "1")
			return next(ctx)
		}
	}
	var Middleware2 = func(next RouteHandler) RouteHandler {
		return func(ctx Context) error {
			ctx.Set("middleware-order", ctx.Value("middleware-order").(string)+"2")
			return next(ctx)
		}
	}
	var Middleware3 = func(next RouteHandler) RouteHandler {
		return func(ctx Context) error {
			ctx.Set("middleware-order", ctx.Value("middleware-order").(string)+"3")
			return next(ctx)
		}
	}
	app := setupAppTest()
	app.Use(Middleware1, Middleware2, Middleware3)
	app.Get("/hello", func(ctx Context) error {
		if ctx.Value("middleware-order").(string) != "123" {
			t.Error("middleware-order should be 123")
		}
		return ctx.RenderJson(map[string]string{"hello": "world"})
	})
	assertRequestStatus(t, app, "GET", "/hello", 200)
}
