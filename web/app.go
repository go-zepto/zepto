package web

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/go-webpack/webpack"
	"github.com/go-zepto/zepto/web/renderer"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
	"os/exec"
	"sync"
)

type MuxHandler func(w http.ResponseWriter, r *http.Request)
type RouteHandler func(ctx Context)

type App struct {
	http.Handler
	opts       Options
	rootRouter *mux.Router
	n          *negroni.Negroni
	tmpl       *renderer.Tmpl
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
	tmpl, err := renderer.NewTmpl("templates", ".html", true)
	if err != nil {
		panic(err)
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
		tmpl:       tmpl,
	}
}

func (app *App) Start() {
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
			Context: context.Background(),
			logger:  app.opts.logger,
			broker:  app.opts.broker,
			res:     res,
			req:     req,
			data:    &sync.Map{},
			tmpl:    app.tmpl,
		}
		routeHandler(&ctx)
	})
	return app
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.rootRouter.ServeHTTP(w, r)
}
