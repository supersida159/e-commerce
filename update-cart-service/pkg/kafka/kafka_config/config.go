package kafkaconfig

import (
	"time"
)

// ServiceID represents different services
type ServiceID string

const (
	OrchestratorService ServiceID = "Orchestrator"
	OrderService        ServiceID = "Order"
	InventoryService    ServiceID = "Inventory"
	CartService         ServiceID = "Cart"
)

// TopicConfig holds the configuration for Kafka topics
type TopicConfig struct {
	Name              string
	NumPartitions     int32
	ReplicationFactor int16
	RetentionTime     time.Duration
	Configs           map[string]string
}

// KafkaTopics defines the topics for saga orchestration
var KafkaTopics = struct {
	// Create order topic - single partition, multiple consumers
	CreateOrder string

	// Update order status topic - single partition for ordered updates
	UpdateOrder string

	// Individual rollback topics for each service
	RollbackOrder     string
	RollbackInventory string
	RollbackCart      string

	// Update rollback status topic
	UpdateRollback string
}{
	CreateOrder:       "saga.create.order",
	UpdateOrder:       "saga.update.order",
	RollbackOrder:     "saga.rollback.order",
	RollbackInventory: "saga.rollback.inventory",
	RollbackCart:      "saga.rollback.cart",
	UpdateRollback:    "saga.update.rollback",
}

func GetTopicConfigs() []TopicConfig {
	sagaRetention := 24 * time.Hour * 2 // 2 days retention

	baseConfig := map[string]string{
		"cleanup.policy":         "delete",
		"retention.bytes":        "536870912", // 512MB
		"message.timestamp.type": "CreateTime",
		"max.message.bytes":      "1048576", // 1MB max message size
	}

	return []TopicConfig{
		{
			Name:              KafkaTopics.CreateOrder,
			NumPartitions:     1, // Single partition for ordered processing
			ReplicationFactor: 1,
			RetentionTime:     sagaRetention,
			Configs:           baseConfig,
		},
		{
			Name:              KafkaTopics.UpdateOrder,
			NumPartitions:     1, // Single partition for ordered updates
			ReplicationFactor: 1,
			RetentionTime:     sagaRetention,
			Configs:           baseConfig,
		},
		{
			Name:              KafkaTopics.RollbackOrder,
			NumPartitions:     1, // One partition per rollback topic
			ReplicationFactor: 1,
			RetentionTime:     sagaRetention,
			Configs:           baseConfig,
		},
		{
			Name:              KafkaTopics.RollbackInventory,
			NumPartitions:     1,
			ReplicationFactor: 1,
			RetentionTime:     sagaRetention,
			Configs:           baseConfig,
		},
		{
			Name:              KafkaTopics.RollbackCart,
			NumPartitions:     1,
			ReplicationFactor: 1,
			RetentionTime:     sagaRetention,
			Configs:           baseConfig,
		},
		{
			Name:              KafkaTopics.UpdateRollback,
			NumPartitions:     1,
			ReplicationFactor: 1,
			RetentionTime:     sagaRetention,
			Configs:           baseConfig,
		},
	}
}

// GetConsumerGroupMapping returns consumer group IDs for each service
func GetConsumerGroupMapping() map[ServiceID]string {
	return map[ServiceID]string{
		OrderService:        "order-service-group",
		InventoryService:    "inventory-service-group",
		CartService:         "cart-service-group",
		OrchestratorService: "orchestrator-service-group",
	}
}

// GetServiceTopicMapping returns which topics each service should listen to
func GetServiceTopicMapping() map[ServiceID][]string {
	return map[ServiceID][]string{
		// Orchestrator listens to update status and sends rollbacks
		OrchestratorService: {
			KafkaTopics.UpdateOrder,
			KafkaTopics.UpdateRollback,
		},
		// Services listen to create order and their specific rollback topics
		OrderService: {
			KafkaTopics.CreateOrder,
			KafkaTopics.RollbackOrder,
		},
		InventoryService: {
			KafkaTopics.CreateOrder,
			KafkaTopics.RollbackInventory,
		},
		CartService: {
			KafkaTopics.CreateOrder,
			KafkaTopics.RollbackCart,
		},
	}
}
