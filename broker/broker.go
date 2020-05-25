package broker

import (
	"context"
	log "github.com/sirupsen/logrus"
	"reflect"
	"sync"
)

type SubscribeHandler func(ctx context.Context, msg *Message)

type InitOptions struct {
	Logger *log.Logger
}

type BrokerProvider interface {
	Init(opts *InitOptions)
	Publish(ctx context.Context, topic string, msg *Message) error
	Subscribe(ctx context.Context, topic string, handler SubscribeHandler) error
	GracefulStop(ctx context.Context) error
}

type Message struct {
	Header map[string]string
	Body   []byte
}

// BrokerWrapper is a struct that wrap the broker provider (gcp, rabbitmq, etc) and handle with message encode/decode
type Broker struct {
	logger *log.Logger
	p      BrokerProvider
	mux    sync.Mutex
}

func NewBroker(logger *log.Logger, provider BrokerProvider) *Broker {
	return &Broker{
		logger: logger,
		p:      provider,
	}
}

func (b *Broker) Init(opts *InitOptions) {
	b.p.Init(opts)
}

// Publish is a call to the broker publish with encoded message
func (b *Broker) Publish(ctx context.Context, topic string, m interface{}) error {
	msg, err := encodeMessage(m)
	if err != nil {
		return err
	}
	return b.p.Publish(ctx, topic, msg)
}

// Publish is a call to the broker publish with encoded message
func (b *Broker) Subscribe(ctx context.Context, topic string, handler interface{}) {
	h := func(ctx context.Context, message *Message) {
		objArg := reflect.TypeOf(handler).In(1)
		if objArg.Kind() != reflect.Ptr {
			b.logger.Errorf("Subscription decode error: %s should be a pointer\n", objArg)
			return
		}
		objType := objArg.Elem()
		msg, err := decodeMessage(message, objType)
		if err != nil {
			b.logger.Errorf("Error. Could not decode your message. Make sure your publisher is using the same struct of your subscription [%s]", objType)
			return
		}
		fn := reflect.ValueOf(handler)
		fn.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(msg)})
	}
	err := b.p.Subscribe(ctx, topic, h)
	if err != nil {
		b.logger.Error(err)
	}
}

// GracefulStop graceful stop all subscriptions
func (b *Broker) GracefulStop(ctx context.Context) error {
	return b.p.GracefulStop(ctx)
}
