package web

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestMiddlewareStack_WithRouter(t *testing.T) {
	var GlobalMiddleware = func(next RouteHandler) RouteHandler {
		return func(ctx Context) error {
			ctx.Set("global-middleware", "ok")
			return next(ctx)
		}
	}
	var ApiMiddleware = func(next RouteHandler) RouteHandler {
		return func(ctx Context) error {
			ctx.Set("api-middleware", "ok")
			return next(ctx)
		}
	}
	app := setupAppTest()
	app.Use(GlobalMiddleware)
	app.Get("/global/hello", func(ctx Context) error {
		if ctx.Value("global-middleware") == nil || ctx.Value("global-middleware").(string) != "ok" {
			t.Error("global-middleware should be present in context")
		}
		if ctx.Value("api-middleware") != nil {
			t.Error("api-middleware should not be present in context")
		}
		return ctx.RenderJson(map[string]string{"hello": "world"})
	})
	// With Router
	r := app.Router("/api")
	r.Use(ApiMiddleware)
	r.Get("/hello", func(ctx Context) error {
		if ctx.Value("global-middleware") == nil || ctx.Value("global-middleware").(string) != "ok" {
			t.Error("global-middleware should be present in context")
		}
		if ctx.Value("api-middleware") == nil || ctx.Value("api-middleware").(string) != "ok" {
			t.Error("api-middleware should be present in context")
		}
		return ctx.RenderJson(map[string]string{"hello": "world"})
	})
	assertRequestStatus(t, app, "GET", "/global/hello", 200)
	assertRequestStatus(t, app, "GET", "/api/hello", 200)
}

func TestMiddlewareStack_RootRouter(t *testing.T) {
	runCount := 0
	var Middleware = func(next RouteHandler) RouteHandler {
		return func(ctx Context) error {
			runCount += 1
			return next(ctx)
		}
	}
	app := setupAppTest()
	app.Use(Middleware)
	app.Get("/hello", func(ctx Context) error {
		return ctx.RenderJson(map[string]string{"hello": "world"})
	})
	assertRequestStatus(t, app, "GET", "/hello", 200)
	assert.Equal(t, 1, runCount)
}

func assertCalledMiddlewares(t *testing.T, path string, requestCount int, expectedCalledMiddlewares []string) {
	calledMiddlewares := make([]string, 0)
	var createMiddleware = func(name string) MiddlewareFunc {
		return func(next RouteHandler) RouteHandler {
			return func(ctx Context) error {
				calledMiddlewares = append(calledMiddlewares, name)
				return next(ctx)
			}
		}
	}

	// App (Root)
	app := setupAppTest()
	app.Use(createMiddleware("root-middleware-1"))
	app.Use(createMiddleware("root-middleware-2"))

	// Router
	router := app.Router("/router")
	router.Use(createMiddleware("router-middleware-1"))
	router.Use(createMiddleware("router-middleware-2"))

	// Root Route [GET /hello]
	app.Get("/hello", func(ctx Context) error {
		return ctx.RenderJson(map[string]string{"hello": "world"})
	})
	// Router Route [GET /router/hello]
	router.Get("/hello", func(ctx Context) error {
		return ctx.RenderJson(map[string]string{"hello": "world"})
	})
	app.Init()
	for n := 0; n < requestCount; n++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", path, nil)
		req.Header.Add("Authorization", "Bearer 123")
		app.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	}
	assert.Len(t, app.rootRouter.middleware.stack, 2)
	assert.Equal(t, expectedCalledMiddlewares, calledMiddlewares)
}

func TestMiddlewareStack_RootRouter_With_Router(t *testing.T) {
	assertCalledMiddlewares(
		t,
		"/hello",
		1,
		[]string{
			"root-middleware-1",
			"root-middleware-2",
		},
	)
	assertCalledMiddlewares(
		t,
		"/router/hello",
		1,
		[]string{
			"root-middleware-1",
			"root-middleware-2",
			"router-middleware-1",
			"router-middleware-2",
		},
	)

	assertCalledMiddlewares(
		t,
		"/router/hello",
		2,
		[]string{
			"root-middleware-1",
			"root-middleware-2",
			"router-middleware-1",
			"router-middleware-2",
			"root-middleware-1",
			"root-middleware-2",
			"router-middleware-1",
			"router-middleware-2",
		},
	)
}
