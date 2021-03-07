package broker

import (
	"context"

	"github.com/go-zepto/zepto/web"
)

type BrokerInstance interface {
	Publish(ctx context.Context, topic string, m interface{}) error
}

type BrokerInstanceDefault struct {
	broker *Broker
}

func (bid *BrokerInstanceDefault) Publish(ctx context.Context, topic string, m interface{}) error {
	return bid.broker.publish(ctx, topic, m)
}

func InstanceFromCtx(ctx web.Context) BrokerInstance {
	i := ctx.PluginInstance("broker")
	brokerInstance, ok := i.(BrokerInstance)
	if !ok {
		return nil
	}
	return brokerInstance
}
