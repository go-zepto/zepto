# Web Application

For those who need an application beyond the zepto core, you can easily use the Zepto web module.


# Project Structure

If you created a web application through the CLI, you will have a project structure similar to this:

```
assets
    js
      app.js
    sass
      app.sass
    entry.json
controllers
    Hello.go
public
    images
        logo.png
    favicon.png
templates
    layout
    pages
main.go
webpack.config.js
```

See in detail this structure:

### assets
Here we put all the javascript and sass that we want to develop. Note: this folder is not public. It is used by the Webpack to generate the build located in `public/build`.

The entrypoint definitions are located at: `entry.json`

### controllers

Here is the definition of the controllers, that is, where are the business rules of the web application.

### public

All public files must be available here

### templates

This is the folder where we put our templates. We use pongo2 as the default template engine, but with Zepto Web it is possible to easily couple another template engine.



# Routing

Routing with Zepto Web is very simple, just do as follows:

```go
    a.Get("/", MyControllerFunc)
    a.Post("/", AnotherControllerFunc)
```

You can also create a subroute:

```go
    v1 := a.Router("/api/v1")
    v1.Get("/", MyControllerFunc)

    v2 := a.Router("/api/v2")
    v2.Get("/", MyControllerFunc)
```

Subroutes also accepts `hosts` params:

```go
    foo := a.Router("/", web.Hosts("foo.go-zepto.com"))
    foo.Get("/", MyControllerFunc)

    bar := a.Router("/", web.Hosts("bar.go-zepto.com"))
    bar.Get("/", MyControllerFunc)

    dynamic := a.Router("/", web.Hosts("{subdomain}.go-zepto.com"))
    dynamic.Get("/", MyControllerFunc)
```




# Controller

A controller is a function like this:

```go
func MyControllerFunc(ctx web.Context) {
	ctx.Render("pages/index")
}
```

The Controller Context (web.Context) has all the request information as well as everything needed to manage the session and render the template.

See the complete definition of the interface:

```go
type Context interface {
	context.Context
    Vars() map[string]string // URL Param variables
	Set(string, interface{}) // Set a value to be used in template
	SetStatus(status int) Context // Set a http status code before render
	Render(template string) // Render a template
	Logger() *log.Logger // Logger instance
	Broker() *broker.Broker // Broker instance
    Session() *Session // Session instance
}
```

# Middlewares

You can easily create middleware with Zepto, as follows:


```go
func LoggedUserMiddleware(next RouteHandler) RouteHandler {
	return func(ctx Context) error {
		token := ctx.Request().Header.Get("Authorization")
		user := User.LoadUserByToken(token)
		ctx.Set("user", user)
		return next(ctx)
	}
}
```

Through the `Use` method you can apply the middleware to all requests:

Example:
```go

    a.Use(LoggedUserMiddleware, AnotherMiddleware)
    a.Get("/", MyControllerFunc)
    a.Get("/me", AnotherControllerFunc)

```

But if you want middleware to be used only on some controllers, you can apply it directly:

```go
    a.Get("/", MyControllerFunc)
    a.Get("/me", LoggedUserMiddleware(AnotherControllerFunc))
```

# Templating

In zepto our main templating engine is pongo2. It has a syntax similar to Django/Jinja2.

Here's an example:

```html
{% extends "../layouts/default.html" %}

{% block title %}
  Hello, Zepto!
{% endblock %}

{% block body %}
  <div class="gopher">
    <img src="/public/images/gopher.png" />
  </div>
  <h1>
    Hello, Zepto!
  </h1>
  <p>
    {{ helloMessage }}
  </p>
{% endblock %}
```

# Render JSON

Just as it is easy to render templates, you can render JSON like this:

```go
type Person struct {
	Name string `json:"name"`
	Email string `json:"email"`
}

func ExampleJsonController(ctx web.Context) error {
	p := &Person{
		Name: "Carlos Strand",
		Email: "carlosstrand@gmail.com",
	}
	return ctx.RenderJson(p)
}
```

# Debug Mode

If your controller panics or returns an error in development mode `ZEPTO_ENV=development`, you can see a custom error page that provides useful information for debugging.

![Zepto Debugger](debugger.png)


# Session

A default session using the CookieStore is ready for use. For production, you need to set the environment variable `SESSION_SECRET`.

Session usage example:

```go
func LoginUserController(ctx web.Context) error {
    // ...
    ctx.Session().Set("user_id", "12345")
    // ...
}

func GetLoggedController(ctx web.Context) error {
    userID := ctx.Session().Get("user_id")
    // ...
}

func LogoutController(ctx web.Context) error {
    userID := ctx.Session().Get("user_id")
    // ...
}

func LogoutController(ctx web.Context) error {
    ctx.Session().Delete("user_id")
    // ...
}

func ClearAllController(ctx web.Context) error {
    ctx.Session().Clear()
    // ...
}

```

# SASS & JS

Developing in Javascript and styling with SASS is easy with Zepto. That's because we have an easy integration with the webpack.

You can create as many entrypoints as you like to optimize the use of resources on your page.

Example:

```json
{
  "home": [
    "js/common.js",
    "js/home.js",
    "sass/home.sass"
  ],
  "product": [
    "js/common.js",
    "js/product.js",
    "sass/home.sass"
  ]
}
```

In development mode `ZEPTO_ENV=development`, your application will hot-reload for every change in JS/SASS.

To include an asset in the template, just do the following:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  {% asset "app.css" %}
  {% asset "vendor.js" %}
  {% asset "app.js" %}
</head>
</html>
```

Thus, according to the environment that is running dev/prod, the asset will be inserted correctly.

In production, a asset will be rendered with a random hash:

```html
<script type="text/javascript" src="/public/build/app-2c659fd82ce9abd9d8e8.js"></script>
```

In development, the asset will be rendered with the webpack-dev-server url:

```html
<script type="text/javascript" src="http://localhost:3808/public/build/vendor.js"></script>
```

# Deployment

To deploy in production just run the make build:

```bash
make build
```

A `build` folder will be generated with the following structure:

```
- build
    app-service
    templates
    public
```

You can use this bundle to run in production.


You can also build with docker:

```bash
docker build -t app-service .
```

Running the docker image:

```bash
docker run -it -p 8000:8000 app-service
```

# Testing

Zepto has a utility to make it easy to write controller tests.

A minimal test example:

```go
package myapp

import (
    "github.com/go-zepto/zepto/web"
    "strings"
	"testing"
)

func MyController(ctx web.Context) error {
    ctx.SetStatus(401)
    return ctx.RenderJson(map[string]string{"hello": "world"})
}

func TestHandlerRequestStatus(t *testing.T) {
	app := web.NewApp()
	zt := web.NewZeptoTest(t, app)
	res, err := zt.TestHandlerRequest(web.TestHandlerRequestOptions{
		Handler: MyController,
		Method:  "POST",
		Target:  "/",
        Body:    strings.NewReader("body data!"),
		InitialSession: map[string]string{
			"abc": "123", // Useful for testing logged user, for example.
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	res.AssertStatusCode(401)
}
```

## Result Assert Methods
- AssertStatusCode(status int)
- AssertBodyEquals(str string)
- AssertBodyContains(str string)
- AssertSessionValue(key string, value string)
- AssertHeaderValue(key string, value string)
