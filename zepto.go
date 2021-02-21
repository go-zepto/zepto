package zepto

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-zepto/zepto/broker"
	"github.com/go-zepto/zepto/logger"
	"github.com/go-zepto/zepto/logger/logrus"
	"github.com/go-zepto/zepto/mailer"
	"github.com/go-zepto/zepto/web"
	"github.com/go-zepto/zepto/web/renderer"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Zepto struct {
	*web.App
	opts       Options
	grpcServer *grpc.Server
	grpcAddr   string
	httpAddr   string
	httpServer *http.Server
	broker     *broker.Broker
	mailer     mailer.Mailer
	logger     logger.Logger
	startedAt  *time.Time
}

func NewZepto(opts ...Option) *Zepto {
	options := newOptions(opts...)
	z := &Zepto{
		opts: options,
	}
	if options.Logger == nil {
		// Logger not set. Using default logger (logrus)
		l := log.New()
		l.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
		l.SetLevel(log.DebugLevel)
		z.logger = logrus.NewLogrus(l)
	} else {
		z.logger = options.Logger
	}
	z.createApp()
	return z
}

func (z *Zepto) SetupGRPC(addr string, fn func(s *grpc.Server)) {
	z.grpcServer = grpc.NewServer()
	z.grpcAddr = addr
	fn(z.grpcServer)
}

func (z *Zepto) SetupHTTP(addr string) {
	srv := &http.Server{
		Addr: addr,
		Handler: &HTTPZeptoHandler{
			z:       z,
			handler: z,
		},
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	z.httpServer = srv
	z.httpAddr = addr
}

func (z *Zepto) createApp() {
	z.App = web.NewApp()
}

func (z *Zepto) SetupBroker(bp broker.BrokerProvider) {
	z.broker = broker.NewBroker(z.logger, bp)
	z.broker.Init(&broker.InitOptions{
		Logger: z.logger,
	})
}

func (z *Zepto) SetupMailer(m mailer.Mailer) {
	z.mailer = m
}

func (z *Zepto) Broker() *broker.Broker {
	return z.broker
}

func (z *Zepto) Logger() logger.Logger {
	return z.logger
}

func (z *Zepto) InitApp() {
	if z.App != nil {
		z.App.Configure(web.ConfigureOptions{
			Broker:         z.broker,
			Logger:         z.logger,
			Mailer:         z.mailer,
			Env:            z.opts.Env,
			TmplEngine:     z.opts.TmplEngine,
			SessionName:    z.opts.SessionName,
			SessionStore:   z.opts.SessionStore,
			WebpackEnabled: z.opts.WebpackEnabled,
		})
		z.Init()
		z.App.StartWebpackServer()
	}
}

func (z *Zepto) InitMailer() {
	var engine renderer.Engine
	if z.App != nil {
		engine = z.App.RendererEngine()
	}
	if z.mailer != nil {
		z.mailer.Init(&mailer.InitOptions{
			RendererEngine: engine,
		})
	}
}

func (z *Zepto) Start() {
	z.InitApp()
	z.InitMailer()
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
			if z.grpcServer != nil {
				z.logger.Info("Stopping HTTP server...")
				z.httpServer.Shutdown(context.Background())
			}
			if z.broker != nil {
				z.logger.Info("Stopping Broker subscriptions...")
				_ = z.broker.GracefulStop(context.Background())
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
			z.Logger().Infof("HTTP server is listening on %s", z.httpAddr)
			z.httpServer.ListenAndServe()
		}()
	}

	errors := make(chan error)
	<-errors
}
