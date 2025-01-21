package consumerlocal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	goredis "github.com/redis/go-redis/v9"
	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/pkg/config"
	kafkaconfig "github.com/supersida159/e-commerce/api-services/pkg/kafka/kafka_config"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
)

const (
	UpdateChannel         = "updateservices"
	RollbackUpdateChannel = "updaterollback"
	RollbackSingleChannel = "rollbacksingle"
)

type SagaConsumer struct {
	consumer              sarama.Consumer
	serviceID             kafkaconfig.ServiceID
	topics                []string
	updateChannels        sync.Map // map[string]chan entities_orders.OrderEvent
	rollbackChannel       chan entities_orders.OrderEvent
	rollbackSingleChannel chan entities_orders.OrderEvent
	ready                 chan bool
	stopChan              chan struct{}
	wg                    sync.WaitGroup
	appCtx                app_context.AppContext
}

func NewSagaConsumer(
	config *config.Schema,
	brokers []string,
	serviceID kafkaconfig.ServiceID,
	appCtx app_context.AppContext,
) (*SagaConsumer, error) {
	// Create a new Sarama configuration
	newSaramaConfig := sarama.NewConfig()

	// Set Kafka version
	kafkaVersion, err := sarama.ParseKafkaVersion(config.Kafka.KafkaVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Kafka version: %w", err)
	}
	newSaramaConfig.Version = kafkaVersion

	// Consumer group rebalance strategy
	newSaramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	// Offset reset policy
	if config.Kafka.ConsumerOffsetReset == "latest" {
		newSaramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	} else {
		newSaramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	// TLS configuration (if enabled)
	if config.Kafka.EnableTLS {
		newSaramaConfig.Net.TLS.Enable = true
		// Add custom TLS configurations if required
	}

	// Timeouts
	newSaramaConfig.Net.DialTimeout = time.Duration(config.Kafka.Timeout) * time.Millisecond
	newSaramaConfig.Net.ReadTimeout = time.Duration(config.Kafka.Timeout) * time.Millisecond
	newSaramaConfig.Net.WriteTimeout = time.Duration(config.Kafka.Timeout) * time.Millisecond

	// Create a new Sarama consumer
	consumer, err := sarama.NewConsumer(brokers, newSaramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	// Fetch the topic mapping for the service ID
	topicMapping := kafkaconfig.GetServiceTopicMapping()
	topics, ok := topicMapping[serviceID]
	if !ok {
		return nil, fmt.Errorf("no topics found for service ID: %s", serviceID)
	}

	// Initialize the SagaConsumer
	return &SagaConsumer{
		consumer:              consumer,
		serviceID:             serviceID,
		topics:                topics,
		updateChannels:        sync.Map{},
		rollbackChannel:       make(chan entities_orders.OrderEvent, 100),
		rollbackSingleChannel: make(chan entities_orders.OrderEvent, 100),
		ready:                 make(chan bool),
		stopChan:              make(chan struct{}),
		appCtx:                appCtx,
	}, nil
}

func (c *SagaConsumer) Start(ctx context.Context) error {
	log.Printf("[%s] Starting consumer for topics: %v", c.serviceID, c.topics)

	for _, topic := range c.topics {
		partitions, err := c.consumer.Partitions(topic)
		if err != nil {
			return fmt.Errorf("failed to get partitions for topic %s: %w", topic, err)
		}

		for _, partition := range partitions {
			c.wg.Add(1)
			go func(topic string, partition int32) {
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
						log.Printf("Context cancelled, stopping consumer for topic %s, partition %d", topic, partition)
						return
					case <-c.stopChan:
						log.Printf("Stopping consumer for topic %s, partition %d", topic, partition)
						return
					}
				}
			}(topic, partition)
		}
	}

	close(c.ready)
	<-c.stopChan
	return nil
}

func (c *SagaConsumer) getOrCreateUpdateChannel(sagaID string) chan entities_orders.OrderEvent {
	ch, _ := c.updateChannels.LoadOrStore(sagaID, make(chan entities_orders.OrderEvent, 1))
	return ch.(chan entities_orders.OrderEvent)
}

func (c *SagaConsumer) cleanupUpdateChannel(sagaID string) {
	if ch, ok := c.updateChannels.LoadAndDelete(sagaID); ok {
		close(ch.(chan entities_orders.OrderEvent))
	}
}

func (c *SagaConsumer) GetEventChannel(sagaID, channelType string) (chan entities_orders.OrderEvent, error) {
	switch channelType {
	case UpdateChannel:
		return c.getOrCreateUpdateChannel(sagaID), nil
	case RollbackUpdateChannel:
		return c.rollbackChannel, nil
	case RollbackSingleChannel:
		return c.rollbackSingleChannel, nil
	default:
		return nil, fmt.Errorf("unknown channel type: %s", channelType)
	}
}

func (c *SagaConsumer) ProcessMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var event entities_orders.OrderEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	log.Printf("[%s] Received event type %s for saga %s", c.serviceID, event.EventType, event.SagaID)

	switch msg.Topic {
	case kafkaconfig.KafkaTopics.UpdateOrder:
		return c.HandleUpdateStatus(ctx, event)
	case kafkaconfig.KafkaTopics.UpdateRollback:
		return c.HandleUpdateRollback(ctx, event)
	}
	return nil
}

func (c *SagaConsumer) HandleUpdateStatus(ctx context.Context, event entities_orders.OrderEvent) error {
	redisClient := c.appCtx.GetCache()
	if redisClient == nil {
		return common.ErrInternalServerError(fmt.Errorf("redis client is nil"))
	}

	lockKey := fmt.Sprintf("saga:%s:lock", event.SagaID)
	if err := redisClient.LockKey(ctx, lockKey); err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	defer redisClient.UnlockKey(ctx, lockKey)

	updateChan := c.getOrCreateUpdateChannel(event.SagaID)

	var existingSaga entities_orders.OrderEvent
	sagaKey := fmt.Sprintf("saga:%s", event.SagaID)
	err := redisClient.Get(sagaKey, &existingSaga)
	if err != nil {
		if err != goredis.Nil {
			return fmt.Errorf("failed to get saga from Redis: %w", err)
		} else {
			if event.Error == "" {
				select {
				case c.rollbackSingleChannel <- event:
				case <-ctx.Done():
					return ctx.Err()
				default:
					log.Printf("Warning: update channel is full for saga %s", event.SagaID)
				}
				return nil
			}
		}
	}
	if err := c.updateServiceStatus(&existingSaga, event); err != nil {
		return err
	}

	if event.Error != "" {
		existingSaga.Error = event.Error
		redisClient.Remove(sagaKey)
		select {
		case updateChan <- existingSaga:
		case <-ctx.Done():
			return ctx.Err()
		default:
			log.Printf("Warning: update channel is full for saga %s", event.SagaID)
		}
	}

	if existingSaga.AllServicesCompleted() {
		successEvent := existingSaga
		successEvent.EventType = string(entities_orders.EventSuccess)

		select {
		case updateChan <- successEvent:
		case <-ctx.Done():
			return ctx.Err()
		default:
			log.Printf("Warning: update channel is full for saga %s", event.SagaID)
		}

		if err := redisClient.Remove(fmt.Sprintf("saga:%s", existingSaga.SagaID)); err != nil {
			return fmt.Errorf("failed to save completed saga state to Redis: %w", err)
		}
		c.cleanupUpdateChannel(event.SagaID)
		return nil
	}

	if err := redisClient.SetWithExpirationPreserve(ctx, sagaKey, existingSaga); err != nil {
		return fmt.Errorf("failed to save updated saga state to Redis: %w", err)
	}

	return nil
}

func (c *SagaConsumer) HandleUpdateRollback(ctx context.Context, event entities_orders.OrderEvent) error {
	if event.Error != "" {
		return nil
	}

	select {
	case c.rollbackChannel <- event:
		log.Printf("Rollback initiated for service: %s with ID: %s", event.CurrentService, event.SagaID)
	case <-ctx.Done():
		return ctx.Err()
	default:
		log.Printf("Warning: rollback channel is full, dropping event for saga %s", event.SagaID)
	}

	return nil
}

func (c *SagaConsumer) updateServiceStatus(existingSaga *entities_orders.OrderEvent, event entities_orders.OrderEvent) error {
	switch event.CurrentService {
	case entities_orders.OrderService:
		return existingSaga.UpdateServiceStatus(
			entities_orders.OrderService,
			event.ServiceStatus.ServiceStates[entities_orders.OrderService].Status,
			event.ServiceStatus.ServiceStates[entities_orders.OrderService].Error,
		)
	case entities_orders.InventoryService:
		return existingSaga.UpdateServiceStatus(
			entities_orders.InventoryService,
			event.ServiceStatus.ServiceStates[entities_orders.InventoryService].Status,
			event.ServiceStatus.ServiceStates[entities_orders.InventoryService].Error,
		)
	case entities_orders.CartService:
		return existingSaga.UpdateServiceStatus(
			entities_orders.CartService,
			event.ServiceStatus.ServiceStates[entities_orders.CartService].Status,
			event.ServiceStatus.ServiceStates[entities_orders.CartService].Error,
		)
	default:
		return fmt.Errorf("unknown service: %s", event.CurrentService)
	}
}

func (c *SagaConsumer) Stop() error {
	close(c.stopChan)
	c.wg.Wait()

	// Cleanup all update channels
	c.updateChannels.Range(func(key, value interface{}) bool {
		close(value.(chan entities_orders.OrderEvent))
		return true
	})

	// Close rollback channels
	close(c.rollbackChannel)
	close(c.rollbackSingleChannel)

	if err := c.consumer.Close(); err != nil {
		return fmt.Errorf("failed to close consumer: %w", err)
	}

	return nil
}

// DeleteUpdateChannel deletes the update channel for a specific saga ID
func (c *SagaConsumer) DeleteUpdateChannel(sagaID string) error {
	if ch, ok := c.updateChannels.LoadAndDelete(sagaID); ok {
		close(ch.(chan entities_orders.OrderEvent))
		return nil
	}
	return fmt.Errorf("update channel not found for saga ID: %s", sagaID)
}
