package gcp

type Options struct {
	ProjectID          string
	TopicPrefix        string
	SubscriptionPrefix string
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	opt := Options{
		SubscriptionPrefix: "sub.",
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// ProjectID is a required Google Pub/Sub project id
func ProjectID(p string) Option {
	return func(o *Options) {
		o.ProjectID = p
	}
}

// TopicPrefix configures a prefix for all topics in subscribers and publishers.
func TopicPrefix(p string) Option {
	return func(o *Options) {
		o.TopicPrefix = p
	}
}

// SubscriptionPrefix configures a prefix for all subs.
func SubscriptionPrefix(p string) Option {
	return func(o *Options) {
		o.SubscriptionPrefix = p
	}
}
