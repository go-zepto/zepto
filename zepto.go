package zepto

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fatih/color"
	"github.com/go-zepto/zepto/config"
	"github.com/go-zepto/zepto/database"
	"github.com/go-zepto/zepto/logger"
	"github.com/go-zepto/zepto/logger/logrus"
	"github.com/go-zepto/zepto/utils"
	"github.com/go-zepto/zepto/web"
	"github.com/go-zepto/zepto/web/renderer/pongo2"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Zepto struct {
	*web.App
	config     config.Config
	opts       Options
	grpcServer *grpc.Server
	grpcAddr   string
	httpAddr   string
	httpServer *http.Server
	logger     logger.Logger
	startedAt  *time.Time
	plugins    map[string]Plugin
	db         *gorm.DB
}

func NewZepto(configs ...config.Config) *Zepto {
	env := utils.GetEnv("ZEPTO_ENV", "development")
	cfg := config.Config{}
	if len(configs) > 0 {
		cfg = configs[0]
	} else {
		configFromFile, err := config.NewConfigFromFile()
		if err != nil {
			panic(err)
		}
		cfg = *configFromFile
	}
	options := Options{
		Name:           cfg.App.Name,
		Version:        cfg.App.Version,
		Env:            env,
		WebpackEnabled: cfg.App.WebpackEnabled,
		SessionName:    cfg.App.Session.Name,
		SessionSecret:  cfg.App.Session.Secret,
		TmplEngine: pongo2.NewPongo2Engine(
			pongo2.TemplateDir("templates"),
			pongo2.Ext(".html"),
			pongo2.AutoReload(env == "development"),
		),
	}
	z := &Zepto{
		opts:    options,
		plugins: make(map[string]Plugin),
		config:  cfg,
	}
	if options.Logger == nil {
		// Logger not set. Using default logger (logrus)
		l := log.New()
		l.SetFormatter(&log.TextFormatter{
			FullTimestamp:    true,
			DisableTimestamp: !cfg.Logger.Timestamp,
			DisableColors:    !cfg.Logger.Colors,
		})
		if !cfg.Logger.Colors {
			color.NoColor = true
		}
		logLevel, err := log.ParseLevel(cfg.Logger.Level)
		if err != nil {
			panic(err)
		}
		l.SetLevel(logLevel)
		z.logger = logrus.NewLogrus(l)
	} else {
		z.logger = options.Logger
	}
	z.createApp()
	z.httpAddr = fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	z.setupDB()
	z.setupHTTP(z.httpAddr)
	return z
}

func (z *Zepto) SetupGRPC(addr string, fn func(s *grpc.Server)) {
	z.grpcServer = grpc.NewServer()
	z.grpcAddr = addr
	fn(z.grpcServer)
}

func (z *Zepto) createDefaultHTTPServer() *http.Server {
	return &http.Server{
		WriteTimeout: time.Duration(z.config.Server.WriteTimeout) * time.Microsecond,
		ReadTimeout:  time.Duration(z.config.Server.ReadTimeout) * time.Millisecond,
	}
}

func (z *Zepto) setupDB() {
	dc := z.config.DB
	if !dc.Enabled {
		return
	}
	conn := database.Connection{
		Adapter:    dc.Adapter,
		Host:       dc.Host,
		Port:       dc.Port,
		Username:   dc.Username,
		Password:   dc.Password,
		Datababase: dc.Database,
	}
	db, err := conn.Open(logger.NewDBLogger(z.logger))
	if err != nil {
		z.logger.Fatal("exiting due to failed database connection")
	}
	z.db = db
}

func (z *Zepto) setupHTTP(addr string) {
	srv := z.createDefaultHTTPServer()
	if z.opts.HTTPServer != nil {
		srv = z.opts.HTTPServer
	}
	srv.Addr = addr
	srv.Handler = &HTTPZeptoHandler{
		z:       z,
		handler: z,
	}
	z.httpServer = srv
	z.httpAddr = addr
}

func (z *Zepto) createApp() {
	z.App = web.NewApp()
}

func (z *Zepto) Logger() logger.Logger {
	return z.logger
}

func (z *Zepto) AddPlugin(plugin Plugin) {
	z.plugins[plugin.Name()] = plugin
	plugin.OnCreated(z)
	z.UsePrepend(plugin.PrependMiddlewares()...)
	z.Use(plugin.AppendMiddlewares()...)
}

func (z *Zepto) Plugins() map[string]Plugin {
	return z.plugins
}

func (z *Zepto) InitApp() {
	for _, p := range z.plugins {
		p.OnZeptoInit(z)
	}
	if z.App != nil {
		opts := web.ConfigureOptions{
			Logger:         z.logger,
			Env:            z.opts.Env,
			TmplEngine:     z.opts.TmplEngine,
			SessionName:    z.opts.SessionName,
			SessionSecret:  z.opts.SessionSecret,
			SessionStore:   z.opts.SessionStore,
			WebpackEnabled: z.opts.WebpackEnabled,
			DB:             z.db,
		}
		pluginInstances := make(map[string]interface{})
		for _, p := range z.plugins {
			pluginInstances[p.Name()] = p.Instance()
		}
		opts.PluginInstances = pluginInstances
		z.App.Configure(opts)
		z.Init()
		z.App.StartWebpackServer()
	}
}

func (z *Zepto) DB() *gorm.DB {
	return z.db
}

func (z *Zepto) Start() {
	z.InitApp()
	now := time.Now()
	z.startedAt = &now
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	z.logger.Infof("Zepto is starting in %s mode...", z.opts.Env)
	go func() {
		select {
		case sig := <-c:
			z.logger.Infof("Got %s signal.", sig)
			if z.grpcServer != nil {
				z.logger.Info("Stopping gRPC server...")
				z.grpcServer.GracefulStop()
			}
			if z.httpServer != nil {
				z.logger.Info("Stopping HTTP server...")
				z.httpServer.Shutdown(context.Background())
			}
			for pluginName, p := range z.plugins {
				z.logger.Infof("Stopping plugin %s...", pluginName)
				p.OnZeptoStop(z)
			}
			os.Exit(0)
		}
	}()

	if z.grpcServer != nil {
		go func() {
			lis, err := net.Listen("tcp", z.grpcAddr)
			if err != nil {
				z.Logger().Error(err)
				os.Exit(1)
			}
			z.Logger().Infof("gRPC server is listening on %s", z.grpcAddr)
			z.grpcServer.Serve(lis)
		}()
	}

	if z.httpServer != nil {
		go func() {
			z.Logger().Infof("HTTP server is listening on http://%s", z.httpAddr)
			z.httpServer.ListenAndServe()
		}()
	}

	for _, p := range z.plugins {
		p.OnZeptoStart(z)
	}

	errors := make(chan error)
	<-errors
}
