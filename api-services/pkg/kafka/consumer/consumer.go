package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/supersida159/e-commerce/api-services/pkg/kafka/saga"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
)

type SagaStep struct {
	ServiceID  string
	Status     string
	RetryCount int
	Timestamp  time.Time
}

// servicesID for consumer
type ServicesID string

const (
	OrderCreated ServicesID = "CREATE_ORDER_SAGA"
	UpdateSaga   ServicesID = "UPDATE_SAGA"
	Rollback     ServicesID = "ROLLBACK"
)

// // OrderSagaState tracks the state of an order throughout the saga
// type OrderSagaState struct {
// 	OrderID     string
// 	Steps       []SagaStep
// 	CurrentStep int
// 	Status      string // Overall saga status: "PENDING", "COMPLETED", "FAILED"
// 	UpdatedAt   time.Time
// }

// EventHandler defines the interface for handling different types of events
type EventHandler interface {
	HandleOrderCreated(ctx context.Context, event *entities_orders.OrderEvent) error
	HandleUpdateSagaTracker(ctx context.Context, event *entities_orders.OrderEvent) error
	HandleRollback(ctx context.Context, event *entities_orders.OrderEvent) error
	PublishToService(ctx context.Context, serviceID string, event *entities_orders.OrderEvent) error
}

// OrderConsumer structure for consuming order-related events
type OrderConsumer struct {
	consumer     sarama.Consumer
	producer     producers.OrderProducer
	eventHandler EventHandler
	SagaStruct   saga.Orchestrator
	topics       map[producers.ServiceID]string
	ready        chan bool
	stopChan     chan struct{}
	wg           sync.WaitGroup
	stateMutex   sync.RWMutex // Protect saga states map
}

// ConsumerConfig holds the configuration for the consumer
type ConsumerConfig struct {
	Brokers      []string
	Topics       map[producers.ServiceID]string
	GroupID      string
	EventHandler EventHandler
}

func NewOrderConsumer(config ConsumerConfig) (*OrderConsumer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	saramaConfig.Producer.Return.Successes = true // Required for sync producer

	// Create consumer
	consumer, err := sarama.NewConsumer(config.Brokers, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	// Create producer
	saramaProducer, err := sarama.NewSyncProducer(config.Brokers, saramaConfig)
	orderProducer := producers.NewOrderProducer(saramaProducer)

	if err != nil {
		consumer.Close()
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &OrderConsumer{
		consumer:     consumer,
		producer:     *orderProducer,
		eventHandler: config.EventHandler,
		topics:       config.Topics,
		ready:        make(chan bool),
		stopChan:     make(chan struct{}),
	}, nil
}

// Start begins consuming messages from all configured topics
func (c *OrderConsumer) Start(ctx context.Context) error {
	// Create list of topics to consume from
	var topics []string
	for _, topic := range c.topics {
		topics = append(topics, topic)
	}

	// Start a consumer for each topic
	for _, topic := range topics {
		c.wg.Add(1)
		go func(topic string) {
			defer c.wg.Done()
			c.consumeTopic(ctx, topic)
		}(topic)
	}

	return nil
}

// Stop gracefully stops the consumer
func (c *OrderConsumer) Stop() error {
	close(c.stopChan)
	c.wg.Wait()
	return c.consumer.Close()
}

func (c *OrderConsumer) consumeTopic(ctx context.Context, topic string) {
	// Create partitionConsumer for the topic
	partitions, err := c.consumer.Partitions(topic)
	if err != nil {
		log.Printf("Failed to get partitions for topic %s: %v", topic, err)
		return
	}

	for _, partition := range partitions {
		pc, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Printf("Failed to start consumer for topic %s partition %d: %v", topic, partition, err)
			continue
		}

		defer pc.Close()

		// Start consuming messages
		for {
			select {
			case msg := <-pc.Messages():
				if err := c.handleMessage(ctx, msg); err != nil {
					log.Printf("Error handling message: %v", err)
				}

			case err := <-pc.Errors():
				log.Printf("Error from consumer: %v", err)

			case <-c.stopChan:
				return

			case <-ctx.Done():
				return
			}
		}
	}
}

func (c *OrderConsumer) handleMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var event entities_orders.OrderEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if msg.Topic == "CREATE_ORDER_SAGA" {
		return c.handleCreateOrderSaga(ctx, &event)
	} else if msg.Topic == "UPDATE_SAGA_TRACKER" {
		return c.handleUpdateSagaTracker(ctx, &event)
	} else if msg.Topic == "ROLLBACK" {
		return c.handleRollback(ctx, &event)
	}
}

func (c *OrderConsumer) handleCreateOrderSaga(ctx context.Context, event *entities_orders.OrderEvent) error {

	return c.SagaStruct.HandleCreateOrderSaga(ctx, *event)
}

func (c *OrderConsumer) handleUpdateSagaTracker(ctx context.Context, event *entities_orders.OrderEvent) error {
	switch event.Metadata["service"] {
	case "orders":
		return c.handleUpdateSagaTrackerOrders(ctx, event)
	}

	return c.updateSagaState(ctx, state)
}

func (c *OrderConsumer) handleRollback(ctx context.Context, event *entities_orders.OrderEvent) error {
	// Reverse through completed steps and rollback each service
	for i := state.CurrentStep - 1; i >= 0; i-- {
		rollbackEvent := &entities_orders.OrderEvent{
			OrderID:   event.OrderID,
			EventType: "ROLLBACK",
			Status:    "PENDING",
		}

		if err := c.eventHandler.PublishToService(ctx, state.Steps[i].ServiceID, rollbackEvent); err != nil {
			log.Printf("Failed to rollback step %d for order %s: %v", i, event.OrderID, err)
		}

		state.Steps[i].Status = "ROLLED_BACK"
	}

	state.Status = "ROLLED_BACK"
	return c.updateSagaState(ctx, state)
}

func mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
