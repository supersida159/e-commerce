package pubsub

import (
	"fmt"
	"time"
)

type Message struct {
	Id       string      `json:"string"`
	channel  Topic       `json:"channel"`
	data     interface{} `json:"data"`
	createAt time.Time   `json:"create_at"`
}

func NewMessage(data interface{}) *Message {
	return &Message{
		Id:       fmt.Sprintf("%d", time.Now().UTC().UnixNano()),
		data:     data,
		createAt: time.Now().UTC(),
	}
}

func (m *Message) String() string {
	return fmt.Sprintf("Message{Id: %s, channel: %s, data: %v, createAt: %s}", m.Id, m.channel, m.data, m.createAt)
}
func (m *Message) Channel() Topic {
	return m.channel
}

func (m *Message) SetChannel(channel Topic) {
	m.channel = channel
}
func (m *Message) Data() interface{} {
	return m.data
}
func (m *Message) CreateAt() time.Time {
	return m.createAt
}
