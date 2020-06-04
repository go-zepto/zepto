package web

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/go-webpack/webpack"
	"github.com/go-zepto/zepto/web/renderer"
	"github.com/go-zepto/zepto/web/renderer/pongo2"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/urfave/negroni"
	"net/http"
	"os"
	"os/exec"
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
		app.opts.logger.Info(m)
	}
	_ = cmd.Wait()
}

func (app *App) setupSession() {
	env := app.opts.env
	if app.opts.sessionStore == nil {
		secret := os.Getenv("SESSION_SECRET")
		if secret == "" {
			if env == "production" {
				app.opts.logger.Fatalf("Missing a required environment variable: SESSION_SECRET")
			} else if env == "development" {
				app.opts.logger.Warn("You will need to setup a SESSION_SECRET in production mode.")
				secret = "development-secret"
			}
		}
		app.opts.sessionStore = sessions.NewCookieStore([]byte(secret))
	}
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
	app := &App{
		opts:       options,
		rootRouter: rootRouter,
		tmplEngine: options.tmplEngine,
	}
	app.setupSession()
	return app
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

func (app *App) getSession(res http.ResponseWriter, req *http.Request) *Session {
	session, _ := app.opts.sessionStore.Get(req, app.opts.sessionName)
	return &Session{
		gSession: session,
		req:      req,
		res:      res,
	}
}

// HandleError recovers from panics gracefully and calls
func (app *App) HandleError(res http.ResponseWriter, req *http.Request, err error) {
	if app.opts.env == "development" {
		res.WriteHeader(500)
		renderer.RenderDevelopmentError(res, req, err)
	} else {
		res.WriteHeader(500)
	}
}

func (app *App) HandleMethod(method string, path string, routeHandler RouteHandler) *App {
	app.rootRouter.HandleFunc(path, func(res http.ResponseWriter, req *http.Request) {
		ctx := NewDefaultContext()
		ctx.logger = app.opts.logger
		ctx.broker = app.opts.broker
		ctx.res = res
		ctx.req = req
		ctx.cookies = &Cookies{
			res: res,
			req: req,
		}
		ctx.session = app.getSession(res, req)
		ctx.tmplEngine = app.tmplEngine
		// Handle Controller Panic
		defer func() {
			if r := recover(); r != nil {
				var e error
				switch t := r.(type) {
				case error:
					e = t
				case string:
					e = fmt.Errorf(t)
				default:
					e = fmt.Errorf(fmt.Sprint(t))
				}
				app.HandleError(res, req, e)
			}
		}()
		err := routeHandler(ctx)
		// Handle Controller Error
		if err != nil {
			app.HandleError(res, req, err)
		}
	}).Methods(method)
	return app
}

func (app *App) Get(path string, routeHandler RouteHandler) *App {
	return app.HandleMethod("GET", path, routeHandler)
}

func (app *App) Post(path string, routeHandler RouteHandler) *App {
	return app.HandleMethod("POST", path, routeHandler)
}

func (app *App) Put(path string, routeHandler RouteHandler) *App {
	return app.HandleMethod("PUT", path, routeHandler)
}

func (app *App) Delete(path string, routeHandler RouteHandler) *App {
	return app.HandleMethod("DELETE", path, routeHandler)
}

func (app *App) Patch(path string, routeHandler RouteHandler) *App {
	return app.HandleMethod("PATCH", path, routeHandler)
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ErrorHandler(app).ServeHTTP(w, r)
}
