package consumerlocal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/supersida159/e-commerce/update-inventory-service/pkg/config"
	kafkaconfig "github.com/supersida159/e-commerce/update-inventory-service/pkg/kafka/kafka_config"
	entities "github.com/supersida159/e-commerce/update-inventory-service/src/model"
)

const (
	CreateOrderChannel = "createorder"
	RollbackChannel    = "rollback"
)

type SagaConsumer struct {
	consumer        sarama.Consumer
	serviceID       kafkaconfig.ServiceID
	topics          []string
	rollbackChannel chan *entities.OrderEvent
	createChannel   chan *entities.OrderEvent
	ready           chan bool
	stopChan        chan struct{}
	wg              sync.WaitGroup
}

func NewSagaConsumer(
	config *config.Schema,
	brokers []string,
	serviceID kafkaconfig.ServiceID,
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
		consumer:        consumer,
		serviceID:       serviceID,
		topics:          topics,
		rollbackChannel: make(chan *entities.OrderEvent, 100),
		createChannel:   make(chan *entities.OrderEvent, 100),
		ready:           make(chan bool),
		stopChan:        make(chan struct{}),
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

func (c *SagaConsumer) GetEventChannel(channelType string) (chan *entities.OrderEvent, error) {
	switch channelType {
	case CreateOrderChannel:
		return c.createChannel, nil
	case RollbackChannel:
		return c.rollbackChannel, nil
	default:
		return nil, fmt.Errorf("unknown channel type: %s", channelType)
	}
}

func (c *SagaConsumer) ProcessMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	event := &entities.OrderEvent{} // Initialize the struct
	if err := json.Unmarshal(msg.Value, event); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	log.Printf("[%s] Received event type %s for saga %s", c.serviceID, event.EventType, event.SagaID)

	switch msg.Topic {
	case kafkaconfig.KafkaTopics.CreateOrder:
		c.createChannel <- event
		fmt.Println("event:", CreateOrderChannel)
		return nil
	case kafkaconfig.KafkaTopics.RollbackCart:
		c.rollbackChannel <- event
		fmt.Println("event:", RollbackChannel)
		return nil
	default:
		return fmt.Errorf("unknown topic: %s", msg.Topic)
	}
}

func (c *SagaConsumer) Stop() error {
	close(c.stopChan)
	c.wg.Wait()

	// Close rollback channels
	close(c.rollbackChannel)
	close(c.createChannel)

	if err := c.consumer.Close(); err != nil {
		return fmt.Errorf("failed to close consumer: %w", err)
	}

	return nil
}
