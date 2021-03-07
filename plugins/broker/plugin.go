package broker

import (
	"context"

	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/web"
)

type SubscriptionsMap map[string]interface{}

type Options struct {
	Provider      BrokerProvider
	Subscriptions SubscriptionsMap
}

func NewBrokerPlugin(opts Options) *BrokerPlugin {
	return &BrokerPlugin{
		broker: NewBroker(opts.Provider),
		subs:   opts.Subscriptions,
	}
}

type BrokerPlugin struct {
	broker   *Broker
	subs     map[string]interface{}
	instance BrokerInstance
}

func (b *BrokerPlugin) Name() string {
	return "broker"
}

func (b *BrokerPlugin) Instance() interface{} {
	return b.instance
}

func (b *BrokerPlugin) PrependMiddlewares() []web.MiddlewareFunc {
	return nil
}

func (b *BrokerPlugin) AppendMiddlewares() []web.MiddlewareFunc {
	return nil
}

func (b *BrokerPlugin) OnCreated(z *zepto.Zepto) {
	b.instance = &BrokerInstanceDefault{
		broker: b.broker,
	}
}

func (b *BrokerPlugin) OnZeptoInit(z *zepto.Zepto) {
	b.broker.Init(&InitOptions{
		Instance: b.instance,
		Logger:   z.Logger(),
	})
}

func (b *BrokerPlugin) OnZeptoStart(z *zepto.Zepto) {
	for topic, handler := range b.subs {
		b.broker.subscribe(context.Background(), topic, handler)
	}
}

func (b *BrokerPlugin) OnZeptoStop(z *zepto.Zepto) {
	b.broker.GracefulStop(context.Background())
}
