package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/pkg/kafka/producers"
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
type ConsumerServicesID string

const (
	OrderCreated ConsumerServicesID = "CREATE_ORDER_SAGA"
	UpdateSaga   ConsumerServicesID = "UPDATE_SAGA"
	Rollback     ConsumerServicesID = "UPDATE_ROLLBACK"
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

// OrderConsumer structure for consuming order-related events
type OrderConsumer struct {
	consumer   sarama.Consumer
	Producer   producers.OrderProducer
	SagaStruct saga.Orchestrator
	topics     map[producers.ServiceID]string
	ready      chan bool
	stopChan   chan struct{}
	wg         sync.WaitGroup
	stateMutex sync.RWMutex // Protect saga states map
}

func NewOrderConsumer(config producers.ConsumerProducerConfig, orderProducer producers.OrderProducer, appCtx app_context.Appcontext) (*OrderConsumer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	saramaConfig.Producer.Return.Successes = true // Required for sync producer

	// Create consumer
	consumer, err := sarama.NewConsumer(config.Brokers, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	// // Create producer
	// saramaProducer, err := sarama.NewSyncProducer(config.Brokers, saramaConfig)
	// orderProducer := producers.NewOrderProducer(saramaProducer)

	if err != nil {
		consumer.Close()
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}
	sagaStruct := saga.NewOrchestrator(&orderProducer, appCtx)

	return &OrderConsumer{
		consumer:   consumer,
		Producer:   orderProducer,
		SagaStruct: *sagaStruct,
		topics:     config.Topics,
		ready:      make(chan bool),
		stopChan:   make(chan struct{}),
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
	// Get all partitions for the topic
	partitions, err := c.consumer.Partitions(topic)
	if err != nil {
		log.Printf("failed to get partitions for topic %s: %v", topic, err)
		return
	}

	for _, partition := range partitions {
		partitionConsumer, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Printf("failed to consume partition %d of topic %s: %v", partition, topic, err)
			continue
		}

		defer partitionConsumer.Close()

		for {
			select {
			case msg := <-partitionConsumer.Messages():
				// Process each message
				c.processMessage(ctx, msg)
			case <-c.stopChan:
				return
			case <-ctx.Done():
				return
			}
		}
	}
}
func (c *OrderConsumer) processMessage(ctx context.Context, msg *sarama.ConsumerMessage) {
	var event entities_orders.OrderEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("failed to unmarshal message: %v", err)
		return
	}

	log.Printf("Received event for saga %s: %s", event.SagaID, event.EventType)

	switch ConsumerServicesID(event.EventType) {
	case OrderCreated:
		if err := c.SagaStruct.StartOrderSaga(ctx, &event.Order); err != nil {
			log.Printf("Failed to start order saga: %v", err)
		}
	case UpdateSaga:
		if err := c.SagaStruct.HandleServiceResponse(ctx, event); err != nil {
			log.Printf("Failed to handle service response: %v", err)
		}
	case Rollback:
		if err := c.SagaStruct.HandleCompensation(ctx, event); err != nil {
			log.Printf("Failed to handle compensation: %v", err)
		}
	default:
		log.Printf("Unhandled event type: %s", event.EventType)
	}
}

// func (c *OrderConsumer) handleCreateOrderSaga(ctx context.Context, event *entities_orders.OrderEvent) error {

// 	return c.SagaStruct.HandleCreateOrderSaga(ctx, *event)
// }

// func (c *OrderConsumer) handleUpdateSagaTracker(ctx context.Context, event *entities_orders.OrderEvent) error {

// 	return c.SagaStruct.HandleUpdateOrderSaga(ctx, *event)
// }

// // // handle sent roll back to services
// // func (c *OrderConsumer) handleSentRollback(ctx context.Context, event *entities_orders.OrderEvent) error {

// // 	return c.SagaStruct.HandleRollbackOrderSaga(ctx, *event)
// // }

// // update roll back
// func (c *OrderConsumer) handleUpdateCompensateSaga(ctx context.Context, event *entities_orders.OrderEvent) error {
// 	return c.SagaStruct.HandleUpdateCompensateOrderSaga(ctx, *event)
// }

// func mustMarshal(v interface{}) []byte {
// 	data, err := json.Marshal(v)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return data
// }
