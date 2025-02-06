package pubsub

import "context"

type Topic string

const (
	UpdateOrdeExpire Topic = "update_order_expire"
	CreateOrder      Topic = "create_order"
)

type KeyEvent string

const (
	OrderExpire KeyEvent = "__keyevent@*__:expired"
)

type Prefix string

const (
	Notification Prefix = "notification-"
	Saga         Prefix = "Saga-"
)

type PubSub interface {
	// Publish publishes a message to the topic.
	Publish(ctx context.Context, channel Topic, data *Message) error
	Subscribe(ctx context.Context, channel Topic) (ch <-chan *Message, close func())
	//Unsubscribe(ctx context.Context, channel Topic)
}
