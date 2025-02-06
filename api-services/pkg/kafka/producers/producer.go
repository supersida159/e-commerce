package producers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/supersida159/e-commerce/api-services/common"
	kafkaconfig "github.com/supersida159/e-commerce/api-services/pkg/kafka/kafka_config"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
	"golang.org/x/sync/errgroup"
)

// OrderProducer handles Kafka message production for order-related events
type OrderProducer struct {
	producer sarama.SyncProducer
	config   *kafkaconfig.KafkaConfig
	mu       sync.Mutex
}

// SendMessageOptions configures message delivery parameters
type SendMessageOptions struct {
	TargetServices []kafkaconfig.ServiceID
	Topic          string
	PartitionKey   string
}

// NewOrderProducer creates a new Kafka producer instance
func NewOrderProducer(kafkaConfig *kafkaconfig.KafkaConfig) (*OrderProducer, error) {
	if err := kafkaConfig.Validate(); err != nil {
		return nil, fmt.Errorf("invalid kafka config: %w", err)
	}

	producer, err := sarama.NewSyncProducer(
		kafkaConfig.Brokers,
		kafkaConfig.SaramaConfig(true),
	)
	if err != nil {
		return nil, common.ErrServiceUnavailable(fmt.Errorf("producer creation failed: %w", err))
	}

	return &OrderProducer{
		producer: producer,
		config:   kafkaConfig,
	}, nil
}

// SendOrderEvent sends an order event to a specific topic
func (p *OrderProducer) SendOrderEvent(ctx context.Context, event entities_orders.OrderEvent, topic string) *common.AppError {
	p.mu.Lock()
	defer p.mu.Unlock()

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return common.ErrInvalidInputData(fmt.Errorf("event marshaling failed: %w", err))
	}

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(event.SagaID),
		Value:     sarama.ByteEncoder(eventJSON),
		Timestamp: time.Now().UTC(),
	}

	if _, _, err := p.producer.SendMessage(msg); err != nil {
		return common.ErrInternalServerError(fmt.Errorf("message send failed: %w", err))
	}

	return nil
}

// SendToServices sends an event to multiple services in parallel
func (p *OrderProducer) SendToServices(ctx context.Context, event entities_orders.OrderEvent, opts SendMessageOptions) *common.AppError {
	g, ctx := errgroup.WithContext(ctx)

	for _, serviceID := range opts.TargetServices {
		serviceID := serviceID // Capture for closure
		g.Go(func() error {
			topics := kafkaconfig.GetServiceTopics(serviceID)
			if len(topics) == 0 {
				return fmt.Errorf("no topics found for service %s", serviceID)
			}

			return p.SendOrderEvent(ctx, event, topics[0])
		})
	}

	if err := g.Wait(); err != nil {
		return common.ErrInternalServerError(fmt.Errorf("parallel send failed: %w", err))
	}
	return nil
}

// SendToMultipleTopics sends an event to multiple topics with error aggregation
func (p *OrderProducer) SendToMultipleTopics(ctx context.Context, event entities_orders.OrderEvent, topics []string) *common.AppError {
	if len(topics) == 0 {
		return common.ErrInvalidRequestParameter(fmt.Errorf("empty topics list"))
	}

	g, ctx := errgroup.WithContext(ctx)

	for _, topic := range topics {
		topic := topic // Capture for closure
		g.Go(func() error {
			return p.SendOrderEvent(ctx, event, topic)
		})
	}

	if err := g.Wait().(*common.AppError); err.RootErr != nil {
		return common.ErrInternalServerError(fmt.Errorf("multi-topic send failed: %w", err))
	}
	return nil
}

// Close gracefully shuts down the producer
func (p *OrderProducer) Close() *common.AppError {
	p.mu.Lock()
	defer p.mu.Unlock()

	if err := p.producer.Close(); err != nil {
		return common.ErrInternalServerError(fmt.Errorf("producer shutdown failed: %w", err))
	}
	log.Println("Kafka producer closed successfully")
	return nil
}

// Convenience methods for common operations
func (p *OrderProducer) SendCreateOrder(ctx context.Context, event entities_orders.OrderEvent) *common.AppError {
	return p.SendOrderEvent(ctx, event, kafkaconfig.KafkaTopics.CreateOrder)
}

func (p *OrderProducer) SendStatusUpdate(ctx context.Context, event entities_orders.OrderEvent) *common.AppError {
	return p.SendOrderEvent(ctx, event, kafkaconfig.KafkaTopics.UpdateOrder)
}

func (p *OrderProducer) SendRollback(ctx context.Context, event entities_orders.OrderEvent, service kafkaconfig.ServiceID) *common.AppError {
	switch service {
	case kafkaconfig.OrderService:
		return p.SendOrderEvent(ctx, event, kafkaconfig.KafkaTopics.RollbackOrder)
	case kafkaconfig.InventoryService:
		return p.SendOrderEvent(ctx, event, kafkaconfig.KafkaTopics.RollbackInventory)
	case kafkaconfig.CartService:
		return p.SendOrderEvent(ctx, event, kafkaconfig.KafkaTopics.RollbackCart)
	default:
		return common.ErrInvalidRequestParameter(fmt.Errorf("invalid service for rollback"))
	}
}
