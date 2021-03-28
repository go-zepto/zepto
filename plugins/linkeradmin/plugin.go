package linkeradmin

import (
	"embed"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/web"
)

//go:embed frontend/webapp/build/*
var webappBuild embed.FS

type Options struct {
	// Admin Config
	Admin *Admin

	// Admin URL Path. Default "/admin"
	Path string

	// Enable/Disable automatic guessing of all resources and fields. Default is true
	AutoGuess *bool
}

type LinkerAdminPlugin struct {
	admin  *Admin
	path   string
	router *web.Router
	schema *Schema
}

func (l *LinkerAdminPlugin) serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)
	path := strings.Replace(req.URL.Path, l.path, "/", 1)
	req, _ = http.NewRequest("GET", path, nil)
	proxy.ServeHTTP(res, req)
}

func NewLinkerAdminPlugin(opts Options) *LinkerAdminPlugin {
	if opts.Path == "" {
		opts.Path = "/admin"
	}
	return &LinkerAdminPlugin{
		admin: opts.Admin,
		path:  opts.Path,
	}
}

func (l *LinkerAdminPlugin) Name() string {
	return "linkeradmin"
}

func (l *LinkerAdminPlugin) Instance() interface{} {
	return nil
}

func (l *LinkerAdminPlugin) PrependMiddlewares() []web.MiddlewareFunc {
	return nil
}

func (l *LinkerAdminPlugin) AppendMiddlewares() []web.MiddlewareFunc {
	return nil
}

func (l *LinkerAdminPlugin) OnCreated(z *zepto.Zepto) {
	l.router = z.Router(l.path)
}

func (l *LinkerAdminPlugin) OnZeptoInit(z *zepto.Zepto) {
	l.schema = &Schema{
		Admin: l.admin,
	}
	l.router.Get("/_schema", func(ctx web.Context) error {
		return ctx.RenderJson(l.schema)
	})

	webappURL := os.Getenv("LINKER_ADMIN_WEBAPP_URL")
	if webappURL != "" {
		proxyCtrl := func(ctx web.Context) error {
			l.serveReverseProxy(webappURL, ctx.Response(), ctx.Request())
			return nil
		}

		l.router.Any("/", proxyCtrl)
		z.Any("/static/{rest:.*}", proxyCtrl)
		return
	}

	l.router.Get("/", func(ctx web.Context) error {
		// Root index.html
		indexHTMLBytes, err := webappBuild.ReadFile("frontend/webapp/build/index.html")
		if err != nil {
			return errors.New("could not load admin")
		}
		indexHTML := string(indexHTMLBytes)
		indexHTML = strings.ReplaceAll(indexHTML, "/admin", l.path)
		ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err = fmt.Fprint(ctx.Response(), indexHTML)
		return err
	})

	l.router.Get("/{rest:.*}", func(ctx web.Context) error {
		if ctx.Request().URL.Path == l.path+"/" {
			return ctx.Redirect(l.path)
		}
		filePath := strings.Replace(ctx.Request().URL.Path, l.path, "", 1)
		req, _ := http.NewRequest("GET", "frontend/webapp/build"+filePath, nil)
		handler := http.FileServer(http.FS(webappBuild))
		handler.ServeHTTP(ctx.Response(), req)
		return nil
	})
}

func (l *LinkerAdminPlugin) OnZeptoStart(z *zepto.Zepto) {}

func (l *LinkerAdminPlugin) OnZeptoStop(z *zepto.Zepto) {}
