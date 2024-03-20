package pubsublocal

import (
	"context"
	"log"
	"sync"

	"github.com/supersida159/e-commerce/pkg/pubsub"
)

var TopicOrderCreated pubsub.Topic = "TopicOrderCreated"

type OrderIdWithQuantity struct {
	FakeID   string
	OrderId  string
	Quantity int
}

type localPubsub struct {
	messagesQueue chan *pubsub.Message
	mapChannels   map[pubsub.Topic][]chan *pubsub.Message
	mapOrders     map[string]chan bool
	locker        *sync.Mutex
}

func NewPubSub() *localPubsub {
	pb := &localPubsub{
		messagesQueue: make(chan *pubsub.Message, 10000),
		mapChannels:   make(map[pubsub.Topic][]chan *pubsub.Message),
		mapOrders:     make(map[string]chan bool),
		locker:        &sync.Mutex{},
	}
	pb.run()
	return pb
}

func (pb *localPubsub) run() {

	go func() {
		for {
			select {
			case mess := <-pb.messagesQueue:
				if subs, ok := pb.mapChannels[mess.Channel()]; ok {
					for _, sub := range subs {
						func(c chan *pubsub.Message) {
							c <- mess
						}(sub)
					}
				}
			}

		}
	}()

}

func (pb *localPubsub) Publish(ctx context.Context, topic pubsub.Topic, data *pubsub.Message) error {
	data.SetChannel(topic)
	go func() {
		pb.messagesQueue <- data
		log.Println("new Publish event", data.Data())
	}()
	return nil
}

func (pb *localPubsub) Subscribe(ctx context.Context, topic pubsub.Topic) (ch <-chan *pubsub.Message, close func()) {
	// Create a new channel for receiving messages
	chs := make(chan *pubsub.Message)

	// Lock the mutex to ensure atomicity in map access
	pb.locker.Lock()

	// Check if the topic already has subscribers
	if val, ok := pb.mapChannels[topic]; ok {
		// If yes, append the new channel to the list of subscribers
		val = append(val, chs)
		pb.mapChannels[topic] = val
	} else {
		// If no subscribers for the topic, create a new slice with the current channel
		pb.mapChannels[topic] = []chan *pubsub.Message{chs}
	}

	// Unlock the mutex
	pb.locker.Unlock()

	// Return the channel for receiving messages and an unsubscribe closure
	return chs, func() {
		// Lock the mutex before modifying the mapChannels
		pb.locker.Lock()

		// Check if the topic has subscribers
		if val, ok := pb.mapChannels[topic]; ok {
			// Iterate over the list of subscribers to find and remove the current channel
			for i, ch := range val {
				if ch == chs {
					// Remove the channel from the slice
					val = append(val[:i], val[i+1:]...)
					pb.mapChannels[topic] = val
					break
				}
			}
		}

		// Unlock the mutex after modifying the mapChannels
		pb.locker.Unlock()
	}
}
