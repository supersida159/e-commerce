package producers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/supersida159/e-commerce/update-inventory-service/common"
	"github.com/supersida159/e-commerce/update-inventory-service/pkg/config"
	kafkaconfig "github.com/supersida159/e-commerce/update-inventory-service/pkg/kafka/kafka_config"
	entities "github.com/supersida159/e-commerce/update-inventory-service/src/model"
)

// OrderProducer structure for saga pattern
type OrderProducer struct {
	producer sarama.SyncProducer
}

// SendMessageOptions contains options for sending messages
type SendMessageOptions struct {
	TargetServices []kafkaconfig.ServiceID
	Topic          string // Specific topic to send to
}

// NewOrderProducer creates a new instance of OrderProducer
func NewOrderProducer(config *config.Schema) (*OrderProducer, error) {
	// Create a new Sarama configuration object
	saramaConfig := sarama.NewConfig()

	// Set producer configurations using values from config.Schema
	saramaConfig.Producer.RequiredAcks = sarama.RequiredAcks(config.Kafka.ProducerRequiredAcks) // Use value from config.Schema
	saramaConfig.Producer.Return.Successes = true                                               // Always return successful messages
	saramaConfig.Producer.Retry.Max = config.Kafka.Retry                                        // Set retry attempts from config.Schema

	// Optional: Set timeout from config.Schema (if provided)
	saramaConfig.Net.DialTimeout = time.Duration(config.Kafka.Timeout) * time.Millisecond
	saramaConfig.Net.ReadTimeout = time.Duration(config.Kafka.Timeout) * time.Millisecond
	saramaConfig.Net.WriteTimeout = time.Duration(config.Kafka.Timeout) * time.Millisecond

	// Configure TLS if enabled in config.Schema
	if config.Kafka.EnableTLS {
		saramaConfig.Net.TLS.Enable = true
		// You can add custom TLS settings here (e.g., cert files, etc.)
	}

	// Create the Kafka producer using the Sarama configuration
	producer, err := sarama.NewSyncProducer(config.Kafka.Brokers, saramaConfig)
	if err != nil {
		// Return a service unavailable error if the producer creation fails
		return nil, common.ErrServiceUnavailable(fmt.Errorf("failed to create producer: %w", err))
	}

	// Return the OrderProducer object with the producer
	return &OrderProducer{
		producer: producer,
	}, nil
}

// SendCreateOrder sends an order creation event to all services
func (p *OrderProducer) SendCreateOrder(event entities.OrderEvent) *common.AppError {
	return p.sendToTopic(event, kafkaconfig.KafkaTopics.CreateOrder)
}

// SendStatusUpdate sends a status update to the orchestrator
func (p *OrderProducer) SendStatusUpdate(event entities.OrderEvent) *common.AppError {
	return p.sendToTopic(event, kafkaconfig.KafkaTopics.UpdateOrder)
}

// SendRollbackOrder sends rollback commands to specified services
func (p *OrderProducer) SendRollbackOrder(event entities.OrderEvent) *common.AppError {
	return p.sendToTopic(event, kafkaconfig.KafkaTopics.RollbackOrder)
}

func (p *OrderProducer) SendRollbackInventory(event entities.OrderEvent) *common.AppError {
	return p.sendToTopic(event, kafkaconfig.KafkaTopics.RollbackInventory)
}

func (p *OrderProducer) SendRollbackCart(event entities.OrderEvent) *common.AppError {
	return p.sendToTopic(event, kafkaconfig.KafkaTopics.RollbackCart)
}

func (p *OrderProducer) SendUpdateRollback(event entities.OrderEvent) *common.AppError {
	return p.sendToTopic(event, kafkaconfig.KafkaTopics.UpdateRollback)
}

// sendToTopic sends an event to a specific topic
func (p *OrderProducer) sendToTopic(event entities.OrderEvent, topic string) *common.AppError {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return common.ErrInvalidInputData(fmt.Errorf("failed to marshal event: %w", err))
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(event.SagaID),
		Value: sarama.StringEncoder(eventJSON),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return common.ErrInternalServerError(fmt.Errorf("failed to send message: %w", err))
	}

	log.Printf("Message sent to topic %s [partition: %d, offset: %d]", topic, partition, offset)
	return nil
}

// SendToMultipleTopics sends an event to multiple topics in parallel
func (p *OrderProducer) SendToMultipleTopics(event entities.OrderEvent, topics []string) *common.AppError {
	if len(topics) == 0 {
		return common.ErrInvalidRequestParameter(fmt.Errorf("no topics provided"))
	}

	var wg sync.WaitGroup
	errorsChan := make(chan *common.AppError, len(topics))

	for _, topic := range topics {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			if err := p.sendToTopic(event, t); err != nil {
				errorsChan <- err
			}
		}(topic)
	}

	wg.Wait()
	close(errorsChan)

	var errors []string
	for err := range errorsChan {
		errors = append(errors, err.MessageEn)
	}

	if len(errors) > 0 {
		return common.ErrInternalServerError(fmt.Errorf("multiple errors: %s", strings.Join(errors, "; ")))
	}

	return nil
}

// SendParallel sends messages to multiple services in parallel
func (p *OrderProducer) SendParallel(event entities.OrderEvent, opts SendMessageOptions) (map[kafkaconfig.ServiceID]*common.AppError, *common.AppError) {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return nil, common.ErrInvalidInputData(fmt.Errorf("failed to marshal event: %w", err))
	}

	// If specific topic is provided
	if opts.Topic != "" {
		msg := &sarama.ProducerMessage{
			Topic: opts.Topic,
			Key:   sarama.StringEncoder(event.SagaID),
			Value: sarama.StringEncoder(eventJSON),
		}

		_, _, err := p.producer.SendMessage(msg)
		if err != nil {
			return nil, common.ErrInternalServerError(fmt.Errorf("failed to send message to topic %s: %w", opts.Topic, err))
		}
		return nil, nil
	}

	// Get topic mapping
	topicMapping := kafkaconfig.GetServiceTopicMapping()

	if len(opts.TargetServices) == 0 {
		for service := range topicMapping {
			opts.TargetServices = append(opts.TargetServices, service)
		}
	}

	var wg sync.WaitGroup
	results := make(map[kafkaconfig.ServiceID]*common.AppError)
	var resultsLock sync.Mutex

	for _, serviceID := range opts.TargetServices {
		wg.Add(1)
		go func(sID kafkaconfig.ServiceID) {
			defer wg.Done()

			topics, exists := topicMapping[sID]
			if !exists || len(topics) == 0 {
				resultsLock.Lock()
				results[sID] = common.ErrResourceNotFound(fmt.Errorf("no topic mapping found for service %s", sID))
				resultsLock.Unlock()
				return
			}

			msg := &sarama.ProducerMessage{
				Topic: topics[0],
				Key:   sarama.StringEncoder(event.SagaID),
				Value: sarama.StringEncoder(eventJSON),
			}

			_, _, err := p.producer.SendMessage(msg)

			resultsLock.Lock()
			if err != nil {
				results[sID] = common.ErrInternalServerError(err)
			}
			resultsLock.Unlock()
		}(serviceID)
	}

	wg.Wait()

	// Check for errors
	hasError := false
	for _, err := range results {
		if err != nil {
			hasError = true
			break
		}
	}

	if hasError {
		return results, common.ErrInternalServerError(fmt.Errorf("some services failed to receive messages"))
	}

	return results, nil
}

// Close closes the producer
func (p *OrderProducer) Close() *common.AppError {
	if err := p.producer.Close(); err != nil {
		return common.ErrInternalServerError(fmt.Errorf("failed to close producer: %w", err))
	}
	return nil
}
