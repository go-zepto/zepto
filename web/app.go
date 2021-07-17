package web

import (
	"bufio"
	"errors"
	"net/http"
	"os"
	"os/exec"
	pathlib "path"

	log "github.com/sirupsen/logrus"

	"github.com/go-webpack/webpack"
	"github.com/go-zepto/zepto/logger"
	"github.com/go-zepto/zepto/utils"
	"github.com/go-zepto/zepto/web/renderer"
	"github.com/go-zepto/zepto/web/renderer/pongo2"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/urfave/negroni"
)

type MuxHandler func(w http.ResponseWriter, r *http.Request)
type RouteHandler func(ctx Context) error
type MiddlewareFunc func(RouteHandler) RouteHandler

type Options struct {
	logger          logger.Logger
	env             string
	webpackEnabled  bool
	tmplEngine      renderer.Engine
	sessionName     string
	sessionStore    sessions.Store
	pluginInstances map[string]interface{}
}

type App struct {
	http.Handler
	opts       Options
	muxRouter  *mux.Router
	n          *negroni.Negroni
	tmplEngine renderer.Engine
	rootRouter *Router
	routers    []*Router
}

type ConfigureOptions struct {
	Logger          logger.Logger
	Env             string
	WebpackEnabled  bool
	TmplEngine      renderer.Engine
	SessionName     string
	SessionStore    sessions.Store
	PluginInstances map[string]interface{}
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

func NewApp() *App {
	muxRouter := mux.NewRouter()
	staticDir := "/public/"
	// Create the static router
	muxRouter.
		PathPrefix(staticDir).
		Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))
	app := &App{
		muxRouter:  muxRouter,
		routers:    make([]*Router, 0),
		rootRouter: NewRouter(""),
	}
	// Configure defaults
	env := utils.GetEnv("ZEPTO_ENV", "development")
	app.Configure(ConfigureOptions{
		Env:            env,
		Logger:         log.New(),
		WebpackEnabled: true,
		SessionName:    "zsid",
		TmplEngine: pongo2.NewPongo2Engine(
			pongo2.TemplateDir("templates"),
			pongo2.Ext(".html"),
			pongo2.AutoReload(env == "development"),
		),
	})
	return app
}

func (app *App) Configure(opts ConfigureOptions) {
	app.opts = Options{
		logger:          opts.Logger,
		env:             opts.Env,
		sessionName:     opts.SessionName,
		sessionStore:    opts.SessionStore,
		tmplEngine:      opts.TmplEngine,
		webpackEnabled:  opts.WebpackEnabled,
		pluginInstances: opts.PluginInstances,
	}
	app.tmplEngine = app.opts.tmplEngine
}

func (app *App) Init() {
	// Setup Session
	app.setupSession()
	// Initialize Root Router Handlers
	app.initRouterHandlers(app.rootRouter)
	// Initialize Router Hanlders
	for _, router := range app.routers {
		app.initRouterHandlers(router)
	}
	// Initialize Template Engine
	app.tmplEngine.Init()
}

func (app *App) StartWebpackServer() {
	dev := app.opts.env == "development"
	if _, err := os.Stat("webpack.config.js"); dev && os.IsNotExist(err) {
		return
	}
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

func (app *App) isDev() bool {
	return app.opts.env == "development"
}

func (app *App) getSession(res http.ResponseWriter, req *http.Request) *Session {
	session, _ := app.opts.sessionStore.Get(req, app.opts.sessionName)
	return &Session{
		gSession: session,
		req:      req,
		res:      res,
	}
}

func (app *App) HandleMethod(methods []string, path string, routeHandler RouteHandler) *App {
	app.rootRouter.HandleMethod(methods, path, routeHandler)
	return app
}

func (app *App) Get(path string, routeHandler RouteHandler) *App {
	return app.HandleMethod([]string{"GET"}, path, routeHandler)
}

func (app *App) Post(path string, routeHandler RouteHandler) *App {
	return app.HandleMethod([]string{"POST"}, path, routeHandler)
}

func (app *App) Put(path string, routeHandler RouteHandler) *App {
	return app.HandleMethod([]string{"PUT"}, path, routeHandler)
}

func (app *App) Delete(path string, routeHandler RouteHandler) *App {
	return app.HandleMethod([]string{"DELETE"}, path, routeHandler)
}

func (app *App) Patch(path string, routeHandler RouteHandler) *App {
	return app.HandleMethod([]string{"PATCH"}, path, routeHandler)
}

func (app *App) Any(path string, routeHandler RouteHandler) *App {
	return app.HandleMethod([]string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}, path, routeHandler)
}

func (app *App) Use(mw ...MiddlewareFunc) {
	app.rootRouter.middleware.Use(mw...)
}

func (app *App) UsePrepend(mw ...MiddlewareFunc) {
	app.rootRouter.middleware.UsePrepend(mw...)
}

func (app *App) Resource(path string, resource Resource) *App {
	id_path := pathlib.Join(path, "/{id}")
	app.Get(path, resource.List)
	app.Get(id_path, resource.Show)
	app.Post(path, resource.Create)
	app.Put(id_path, resource.Update)
	app.Delete(id_path, resource.Destroy)
	return app
}

//func (app *App) Use(mw ...MiddlewareFunc)

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.muxRouter.ServeHTTP(w, r)
}

func (app *App) RendererEngine() renderer.Engine {
	return app.tmplEngine
}
