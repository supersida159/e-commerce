package pubsub

import "context"

type Topic string
type PubSub interface {
	// Publish publishes a message to the topic.
	Publish(ctx context.Context, channel Topic, data *Message) error
	Subscribe(ctx context.Context, channel Topic) (ch <-chan *Message, close func())
	//Unsubscribe(ctx context.Context, channel Topic)
}
