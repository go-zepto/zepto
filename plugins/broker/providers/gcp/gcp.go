package gcp

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/dustin/go-humanize"
	"github.com/go-zepto/zepto/logger"
	"github.com/go-zepto/zepto/plugins/broker"
	goption "google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Broker struct {
	logger logger.Logger
	opts   Options
	client *pubsub.Client
	// subscriptions is a map of subId -> subscription
	subscriptions map[string]*subscription
	instance      broker.BrokerInstance
}

type subscription struct {
	b       *Broker
	id      string
	topicId string
	sub     *pubsub.Subscription
	handler broker.SubscribeHandler
	exit    chan bool
}

func NewBroker(opts ...Option) broker.BrokerProvider {
	options := newOptions(opts...)
	gcpOpts := make([]goption.ClientOption, 0)
	if options.CredentialsJSON != nil {
		gcpOpts = append(gcpOpts, goption.WithCredentialsJSON(options.CredentialsJSON))
	}
	c, err := pubsub.NewClient(context.Background(), options.ProjectID, gcpOpts...)
	if err != nil {
		panic(err)
	}
	return &Broker{
		opts:          options,
		client:        c,
		subscriptions: make(map[string]*subscription),
	}
}

func (b *Broker) Init(opts *broker.InitOptions) {
	b.logger = opts.Logger
	b.instance = opts.Instance
}

func (b *Broker) getOrCreateTopic(ctx context.Context, topicId string) (*pubsub.Topic, error) {
	t := b.client.Topic(topicId)
	exists, err := t.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		return t, nil
	}
	return b.client.CreateTopic(ctx, topicId)
}

func (b *Broker) getOrCreateSubscription(ctx context.Context, subId string, topicId string) (*pubsub.Subscription, error) {
	t, err := b.getOrCreateTopic(ctx, topicId)
	if err != nil {
		return nil, err
	}
	sub := b.client.Subscription(subId)
	exists, err := sub.Exists(ctx)
	if exists {
		return sub, nil
	}
	return b.client.CreateSubscription(ctx, subId, pubsub.SubscriptionConfig{
		Topic: t,
	})
}

func (b Broker) Publish(ctx context.Context, topicId string, msg *broker.Message) error {
	topicId = b.opts.TopicPrefix + topicId
	t := b.client.Topic(topicId)

	m := &pubsub.Message{
		Data:       msg.Body,
		Attributes: msg.Header,
	}

	size := len(msg.Body)
	b.logger.Debugf("Publishing message to topic %s (size=%s)", topicId, humanize.Bytes(uint64(size)))
	pr := t.Publish(ctx, m)
	if _, err := pr.Get(ctx); err != nil {
		// create Topic if not exists
		if status.Code(err) == codes.NotFound {

			b.logger.Infof("Topic %s does not exist. Creating...", topicId)
			if t, err = b.client.CreateTopic(ctx, topicId); err == nil {
				_, err = t.Publish(ctx, m).Get(ctx)
			}
		}
	}
	return nil
}

func (b *Broker) Subscribe(ctx context.Context, topicId string, handler broker.SubscribeHandler) error {
	topicId = b.opts.TopicPrefix + topicId
	subID := b.opts.SubscriptionPrefix + topicId
	sub, err := b.getOrCreateSubscription(ctx, subID, topicId)
	if err != nil {
		return err
	}
	s := &subscription{
		b:       b,
		id:      subID,
		topicId: topicId,
		sub:     sub,
		handler: handler,
		exit:    make(chan bool),
	}
	b.logger.Infof("Subscribing to topic %s...", topicId)
	b.subscriptions[subID] = s
	go s.Run()
	return nil
}

func (s *subscription) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	for {
		select {
		case <-s.exit:
			cancel()
			return
		default:
			sub, err := s.b.getOrCreateSubscription(ctx, s.id, s.topicId)
			if err != nil {
				continue
			}
			if err := sub.Receive(ctx, func(ctx context.Context, pm *pubsub.Message) {
				m := &broker.Message{
					Header: pm.Attributes,
					Body:   pm.Data,
				}
				go s.handler(&broker.DefaultSubscriptionContext{
					Context:        ctx,
					BrokerInstance: s.b.instance,
				}, m)
				pm.Ack()
			}); err != nil {
				s.b.logger.Error(err)
				time.Sleep(time.Second)
				continue
			}
		}
	}
}

func (s *subscription) Unsubscribe(ctx context.Context) error {
	s.b.logger.Infof("Unsubscribing %s", s.id)
	return s.sub.Delete(ctx)
}

func (b *Broker) GracefulStop(ctx context.Context) error {
	// TODO: Maybe we can use goroutines to stop all subs
	for _, s := range b.subscriptions {
		close(s.exit)
		s.Unsubscribe(ctx)
	}
	return nil
}
