package zepto

import (
	"context"
	"github.com/carlosstrand/zepto/broker"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
)

type Zepto struct {
	opts       Options
	grpcServer *grpc.Server
	grpcAddr   string
	broker     *broker.Broker
	logger     *log.Logger
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
	return z.logger.WithFields(log.Fields{
		"ts":    "uhuu",
		"other": "I also should be logged always",
	}).Logger
}

func (z *Zepto) Start() error {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	go func() {
		select {
		case sig := <-c:
			z.logger.Infof("Got %s signal.", sig)
			z.logger.Info("Stopping gRPC server...")
			z.grpcServer.GracefulStop()
			z.logger.Info("Stopping Broker subscriptions...")
			_ = z.broker.GracefulStop(context.Background())
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
	return nil
}
