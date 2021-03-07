package testutils

import (
	"context"

	"github.com/go-zepto/zepto/plugins/broker"
)

type PublishCallArgs struct {
	Ctx   context.Context
	Topic string
	Msg   *broker.Message
}

type SubscribeCallArgs struct {
	Ctx     context.Context
	Topic   string
	Handler broker.SubscribeHandler
}

type BrokerProviderMock struct {
	InitCalled        bool
	PublishCalled     bool
	PublishCallArgs   *PublishCallArgs
	SubscribeCalled   bool
	SubscribeCallArgs *SubscribeCallArgs
}

func (p *BrokerProviderMock) Init(opts *broker.InitOptions) {
	p.InitCalled = true
}

func (p *BrokerProviderMock) Publish(ctx context.Context, topic string, msg *broker.Message) error {
	p.PublishCalled = true
	p.PublishCallArgs = &PublishCallArgs{
		Ctx:   ctx,
		Topic: topic,
		Msg:   msg,
	}
	return nil
}

func (p *BrokerProviderMock) Subscribe(ctx context.Context, topic string, handler broker.SubscribeHandler) error {
	p.SubscribeCalled = true
	p.SubscribeCallArgs = &SubscribeCallArgs{
		Ctx:     ctx,
		Topic:   topic,
		Handler: handler,
	}
	return nil
}

func (p *BrokerProviderMock) GracefulStop(ctx context.Context) error {
	return nil
}
