package consumerlocal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	kafkaconfig "github.com/supersida159/e-commerce/api-services/pkg/kafka/kafka_config"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
)

const (
	UpdateChannel         = "update"
	RollbackUpdateChannel = "rollback_update"
	RollbackSingleChannel = "rollback_single"
	channelBufferSize     = 1000
	channelTimeout        = 100 * time.Millisecond
	ttl                   = 5 * time.Minute
)

type SagaConsumer struct {
	consumer  sarama.Consumer
	serviceID kafkaconfig.ServiceID
	topics    []string

	Mu             sync.RWMutex
	updateChannels map[string]chan entities_orders.OrderEvent
	SagaStates     map[string]*entities_orders.OrderEvent

	rollbackUpdateChan chan entities_orders.OrderEvent
	rollbackSingleChan chan entities_orders.OrderEvent
	shutdownChan       chan struct{}

	wg       sync.WaitGroup
	stopOnce sync.Once
}

func NewSagaConsumer(
	config *kafkaconfig.KafkaConfig,
	brokers []string,
	serviceID kafkaconfig.ServiceID,
) (*SagaConsumer, error) {
	consumer, err := sarama.NewConsumer(brokers, config.SaramaConfig(false))
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &SagaConsumer{
		consumer:           consumer,
		serviceID:          serviceID,
		topics:             kafkaconfig.GetServiceTopics(serviceID),
		updateChannels:     make(map[string]chan entities_orders.OrderEvent),
		SagaStates:         make(map[string]*entities_orders.OrderEvent),
		rollbackUpdateChan: make(chan entities_orders.OrderEvent, channelBufferSize),
		rollbackSingleChan: make(chan entities_orders.OrderEvent, channelBufferSize),
		shutdownChan:       make(chan struct{}),
	}, nil
}

func (c *SagaConsumer) GetEventChannel(sagaID string, channelType string) (chan entities_orders.OrderEvent, error) {
	c.Mu.RLock()
	defer c.Mu.RUnlock()

	switch channelType {
	case UpdateChannel:
		return c.getUpdateChannel(sagaID), nil
	case RollbackUpdateChannel:
		return c.rollbackUpdateChan, nil
	case RollbackSingleChannel:
		return c.rollbackSingleChan, nil
	default:
		return nil, fmt.Errorf("unknown channel type: %s", channelType)
	}
}

func (c *SagaConsumer) getUpdateChannel(sagaID string) chan entities_orders.OrderEvent {
	if ch, exists := c.updateChannels[sagaID]; exists {
		return ch
	}
	return c.rollbackSingleChan
}

func (c *SagaConsumer) CreateUpdateChannel(sagaID string) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if _, exists := c.updateChannels[sagaID]; exists {
		return fmt.Errorf("update channel for saga %s already exists", sagaID)
	}

	c.updateChannels[sagaID] = make(chan entities_orders.OrderEvent, channelBufferSize)
	return nil
}

func (c *SagaConsumer) cleanupChannel(sagaID string) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if ch, exists := c.updateChannels[sagaID]; exists {
		close(ch)
		delete(c.updateChannels, sagaID)
	}
	delete(c.SagaStates, sagaID)
}

func (c *SagaConsumer) routeToChannel(
	ctx context.Context,
	sagaID string,
	channelType string,
	event entities_orders.OrderEvent,
) error {
	ch, err := c.GetEventChannel(sagaID, channelType)
	if err != nil {
		return err
	}

	if channelType == RollbackSingleChannel {
		c.Mu.RLock()
		_, updateExists := c.updateChannels[sagaID]
		c.Mu.RUnlock()

		if !updateExists {
			log.Printf("Skipping rollback for saga %s - no active update channel", sagaID)
			return nil
		}
	}

	select {
	case ch <- event:
		if channelType == RollbackUpdateChannel {
			log.Printf("Logged rollback update for saga %s: %v", sagaID, event)
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(channelTimeout):
		return fmt.Errorf("channel buffer overflow for channel type %s", channelType)
	}
}

func (c *SagaConsumer) UpdateSagaState(sagaID string, event *entities_orders.OrderEvent) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.SagaStates[sagaID] = event
}

func (c *SagaConsumer) GetSagaState(sagaID string) (*entities_orders.OrderEvent, bool) {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	state, exists := c.SagaStates[sagaID]
	return state, exists
}

func (c *SagaConsumer) Start(ctx context.Context) error {
	log.Printf("[%s] Starting consumer for topics: %v", c.serviceID, c.topics)

	// Start the TTL-based cleanup routine
	go c.StartCleanupRoutine(ctx, ttl)

	for _, topic := range c.topics {
		partitions, err := c.consumer.Partitions(topic)
		if err != nil {
			return fmt.Errorf("failed to get partitions for topic %s: %w", topic, err)
		}
		for _, partition := range partitions {
			c.wg.Add(1)
			go c.consumePartition(ctx, topic, partition)
		}
	}

	<-c.shutdownChan
	return nil
}

func (c *SagaConsumer) consumePartition(ctx context.Context, topic string, partition int32) {
	defer c.wg.Done()

	partitionConsumer, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Failed to create partition consumer for topic %s, partition %d: %v", topic, partition, err)
		return
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Printf("Failed to close partition consumer: %v", err)
		}
	}()

	log.Printf("[%s] Started consuming topic: %s, partition: %d", c.serviceID, topic, partition)

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			if err := c.ProcessMessage(ctx, msg); err != nil {
				log.Printf("Error processing message: %v", err)
			}
		case err := <-partitionConsumer.Errors():
			log.Printf("Error from partition consumer: %v", err)
		case <-ctx.Done():
			return
		case <-c.shutdownChan:
			return
		}
	}
}

func (c *SagaConsumer) Stop() {
	c.stopOnce.Do(func() {
		close(c.shutdownChan)
		c.wg.Wait()

		c.Mu.Lock()
		for sagaID, ch := range c.updateChannels {
			close(ch)
			delete(c.updateChannels, sagaID)
		}
		c.SagaStates = make(map[string]*entities_orders.OrderEvent)
		c.Mu.Unlock()

		close(c.rollbackUpdateChan)
		close(c.rollbackSingleChan)

		if err := c.consumer.Close(); err != nil {
			log.Printf("Error closing consumer: %v", err)
		}
	})
}

func (c *SagaConsumer) ProcessMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	// Initialize the event variable properly
	event := &entities_orders.OrderEvent{}
	if err := json.Unmarshal(msg.Value, event); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	log.Printf("[%s] Received event type %s for saga %s", c.serviceID, event.EventType, event.SagaID)

	switch msg.Topic {
	case kafkaconfig.KafkaTopics.UpdateOrder:
		return c.HandleUpdateStatus(ctx, event)
	case kafkaconfig.KafkaTopics.UpdateRollback:
		return c.HandleUpdateRollback(ctx, event)
	default:
		return fmt.Errorf("unknown topic: %s", msg.Topic)
	}
}
func (c *SagaConsumer) HandleUpdateStatus(ctx context.Context, event *entities_orders.OrderEvent) error {
	if c.SagaStates[event.SagaID].ServiceStatus.Status != entities_orders.ServicePending &&
		c.SagaStates[event.SagaID].ServiceStatus.Status != entities_orders.ServiceProcessing {
		return c.routeToChannel(ctx, event.SagaID, RollbackSingleChannel, *event)
	} else {
		return c.routeToChannel(ctx, event.SagaID, UpdateChannel, *event)
	}
}

func (c *SagaConsumer) HandleUpdateRollback(ctx context.Context, event *entities_orders.OrderEvent) error {
	return c.routeToChannel(ctx, event.SagaID, RollbackUpdateChannel, *event)
}

func (c *SagaConsumer) CreateNewUpdateChannel(sagaID string) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if _, exists := c.updateChannels[sagaID]; exists {
		return fmt.Errorf("update channel for saga %s already exists", sagaID)
	}

	c.updateChannels[sagaID] = make(chan entities_orders.OrderEvent, channelBufferSize)
	return nil
}
func (c *SagaConsumer) DeleteUpdateChannel(sagaID string) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if ch, exists := c.updateChannels[sagaID]; exists {
		close(ch)
		delete(c.updateChannels, sagaID)
	}
	delete(c.SagaStates, sagaID)
}

func (c *SagaConsumer) GetRollbackUpdateChannel() chan entities_orders.OrderEvent {
	return c.rollbackUpdateChan
}

func (c *SagaConsumer) GetRollbackSingleChannel() chan entities_orders.OrderEvent {
	return c.rollbackSingleChan
}

// StartCleanupRoutine periodically cleans up stale saga states based on TTL
func (c *SagaConsumer) StartCleanupRoutine(ctx context.Context, ttl time.Duration) {
	ticker := time.NewTicker(ttl)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.cleanupStaleSagaStates(ttl)
		}
	}
}

// cleanupStaleSagaStates removes saga states older than the specified TTL
func (c *SagaConsumer) cleanupStaleSagaStates(ttl time.Duration) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	now := time.Now()
	for sagaID, sagaState := range c.SagaStates {
		if now.Sub(sagaState.UpdatedAt) > ttl {
			delete(c.SagaStates, sagaID)
			log.Printf("Deleted stale saga state for saga ID: %s", sagaID)
		}
	}
}
