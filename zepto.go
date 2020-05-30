package zepto

import (
	"context"
	"github.com/go-zepto/zepto/broker"
	"github.com/go-zepto/zepto/web"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Zepto struct {
	opts       Options
	grpcServer *grpc.Server
	grpcAddr   string
	httpAddr   string
	httpServer *http.Server
	broker     *broker.Broker
	logger     *log.Logger
	startedAt  *time.Time
	webApp     *web.App
}

func NewZepto(opts ...Option) *Zepto {
	options := newOptions(opts...)
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetLevel(log.DebugLevel)
	return &Zepto{
		opts:   options,
		logger: logger,
	}
}

func (z *Zepto) SetupGRPC(addr string, fn func(s *grpc.Server)) {
	z.grpcServer = grpc.NewServer()
	z.grpcAddr = addr
	fn(z.grpcServer)
}

func (z *Zepto) SetupHTTP(addr string, handler http.Handler) {
	srv := &http.Server{
		Addr: addr,
		Handler: &HTTPZeptoHandler{
			z:       z,
			handler: handler,
		},
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	z.httpServer = srv
	z.httpAddr = addr
}

func (z *Zepto) NewWeb() *web.App {
	z.webApp = web.NewApp(
		web.Broker(z.broker),
		web.Logger(z.logger),
		web.Env(z.opts.Env),
	)
	return z.webApp
}

func (z *Zepto) SetupBroker(bp broker.BrokerProvider) {
	z.broker = broker.NewBroker(z.logger, bp)
	z.broker.Init(&broker.InitOptions{
		Logger: z.logger,
	})
}

func (z *Zepto) Broker() *broker.Broker {
	return z.broker
}

func (z *Zepto) Logger() *log.Logger {
	return z.logger
}

func (z *Zepto) Start() {
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

	if z.webApp != nil {
		z.webApp.Start()
	}

	errors := make(chan error)
	<-errors
}
