package producer

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
)

// ServiceID represents different services
type ServiceID string

const (
	SagaCentral      ServiceID = "Saga"
	OrderService     ServiceID = "order"
	InventoryService ServiceID = "inventory"
	CartService      ServiceID = "cart"
)

// SendResult stores the result of sending a message to a service
type SendResult struct {
	ServiceID ServiceID
	Error     error
}

// OrderProducer structure with enhanced functionality
type OrderProducer struct {
	producer sarama.SyncProducer
	topics   struct {
		OrderTopic     string
		InventoryTopic string
		CartTopic      string
	}
}

// SendMessageOptions contains options for sending messages
type SendMessageOptions struct {
	// TargetServices specifies which services to send to
	// If nil, sends to all services
	TargetServices []ServiceID
}

// NewOrderProducer creates a new instance of OrderProducer with the given Kafka producer and topics.
func NewOrderProducer(producer sarama.SyncProducer) *OrderProducer {
	return &OrderProducer{
		producer: producer,
		topics: struct {
			OrderTopic     string
			InventoryTopic string
			CartTopic      string
		}{
			OrderTopic:     string(entities_orders.ServiceInit),
			InventoryTopic: string(entities_orders.ServiceInit),
			CartTopic:      string(entities_orders.ServiceInit),
		},
	}
}

func (p *OrderProducer) SendMessages(event entities_orders.OrderEvent, opts *SendMessageOptions) (map[ServiceID]error, error) {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event: %w", err)
	}

	// Define all available services
	serviceTopics := map[ServiceID]string{
		OrderService:     p.topics.OrderTopic,
		InventoryService: p.topics.InventoryTopic,
		CartService:      p.topics.CartTopic,
	}

	// If no specific services are specified, send to all
	targetServices := []ServiceID{OrderService, InventoryService, CartService}
	if opts != nil && len(opts.TargetServices) > 0 {
		targetServices = opts.TargetServices
	}

	// Create messages only for target services
	messages := make(map[ServiceID]*sarama.ProducerMessage)
	for _, serviceID := range targetServices {
		if topic, exists := serviceTopics[serviceID]; exists {
			messages[serviceID] = &sarama.ProducerMessage{
				Topic: topic,
				Key:   sarama.StringEncoder(event.BusinessID),
				Value: sarama.StringEncoder(eventJSON),
			}
		}
	}

	// Send messages concurrently and collect results
	var wg sync.WaitGroup
	results := make(map[ServiceID]error)
	var resultsLock sync.Mutex

	for serviceID, msg := range messages {
		wg.Add(1)
		go func(sID ServiceID, message *sarama.ProducerMessage) {
			defer wg.Done()

			_, _, err := p.producer.SendMessage(message)

			resultsLock.Lock()
			results[sID] = err
			resultsLock.Unlock()
		}(serviceID, msg)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Check if any service failed
	var hasError bool
	for _, err := range results {
		if err != nil {
			hasError = true
			break
		}
	}

	if hasError {
		return results, fmt.Errorf("some services failed to receive messages")
	}

	return results, nil
}

// Example usage functions:

func (p *OrderProducer) SendToAllServices(event entities_orders.OrderEvent) (map[ServiceID]error, error) {
	return p.SendMessages(event, nil)
}

func (p *OrderProducer) SendRollbackToSuccessfulServices(event entities_orders.OrderEvent, successfulServices []ServiceID) (map[ServiceID]error, error) {
	return p.SendMessages(event, &SendMessageOptions{
		TargetServices: successfulServices,
	})
}
