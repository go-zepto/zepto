package web

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/go-webpack/webpack"
	"github.com/go-zepto/zepto/web/renderer"
	"github.com/go-zepto/zepto/web/renderer/pongo2"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
	"os/exec"
	"sync"
)

type MuxHandler func(w http.ResponseWriter, r *http.Request)
type RouteHandler func(ctx Context) error

type App struct {
	http.Handler
	opts       Options
	rootRouter *mux.Router
	n          *negroni.Negroni
	tmplEngine renderer.Engine
}

func (app *App) startWebpack() {
	cmd := exec.Command("npm", "run", "start")
	stdout, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	err := cmd.Start()
	if errors.Is(err, exec.ErrNotFound) {
		app.opts.logger.Errorf("node/npm binaries not found. Please make sure they are installed.")
	}
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	_ = cmd.Wait()
}

func NewApp(opts ...Option) *App {
	options := newOptions(opts...)
	if options.tmplEngine == nil {
		// Use pongo2 as default template engine
		options.tmplEngine = pongo2.NewPongo2Engine(
			pongo2.TemplateDir("templates"),
			pongo2.Ext(".html"),
			pongo2.AutoReload(options.env == "development"),
		)
	}
	rootRouter := mux.NewRouter()
	staticDir := "/public/"
	// Create the static router
	rootRouter.
		PathPrefix(staticDir).
		Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))
	return &App{
		opts:       options,
		rootRouter: rootRouter,
		tmplEngine: options.tmplEngine,
	}
}

func (app *App) Start() {
	err := app.tmplEngine.Init()
	if err != nil {
		panic(err)
	}
	dev := app.opts.env == "development"
	if app.opts.webpackEnabled {
		webpack.FsPath = "./public/build"
		webpack.WebPath = "build"
		webpack.Verbose = true
		webpack.Init(dev)
		if dev {
			go func() {
				app.startWebpack()
			}()
		}
	}
}

func (app *App) GET(path string, routeHandler RouteHandler) *App {
	app.rootRouter.HandleFunc(path, func(res http.ResponseWriter, req *http.Request) {
		ctx := DefaultContext{
			Context:    context.Background(),
			logger:     app.opts.logger,
			broker:     app.opts.broker,
			res:        res,
			req:        req,
			data:       &sync.Map{},
			status:     200,
			tmplEngine: app.tmplEngine,
		}
		err := routeHandler(&ctx)
		if err != nil {
			if app.opts.env == "development" {
				res.WriteHeader(500)
				renderer.RenderDevelopmentError(res, req, err)
			} else {
				res.WriteHeader(500)
			}
		}
	})
	return app
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ErrorHandler(app.rootRouter, app.opts.env).ServeHTTP(w, r)
}
