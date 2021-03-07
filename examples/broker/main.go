package main

import (
	"fmt"

	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/plugins/broker"
	"github.com/go-zepto/zepto/plugins/broker/gcp"
	"github.com/go-zepto/zepto/web"
)

type MyMessage struct {
	Text string `json:"text"`
}

func MyMessageSubscriptionHandler(ctx broker.SubscriptionContext, myMessage *MyMessage) {
	fmt.Printf("Message Received: %s \n", myMessage.Text)
	ctx.Broker().Publish(ctx, "app.example.my_message_received", &MyMessage{Text: "Received"})
}

func main() {
	z := zepto.NewZepto(
		zepto.Name("books-api"),
		zepto.Version("0.0.1"),
	)

	z.AddPlugin(broker.NewBrokerPlugin(broker.BrokerPluginOptions{
		Provider: gcp.NewBroker(gcp.ProjectID("YOUR_GOOGLE_PROJECT_ID")),
		Subscriptions: broker.SubscriptionsMap{
			"app.example.my_message": MyMessageSubscriptionHandler,
		},
	}))

	z.Get("/publish", func(ctx web.Context) error {
		b := broker.InstanceFromCtx(ctx)
		b.Publish(ctx, "app.example.my_message", &MyMessage{
			Text: "Hello World",
		})
		return ctx.RenderJson(map[string]string{
			"satus": "published",
		})
	})

	z.SetupHTTP("0.0.0.0:8000")
	z.Start()
}
