package kafkaconfig

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
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

// InitKafkaTopics creates the Kafka topics based on the configurations
func InitKafkaTopics(brokerList []string) error {
	// Create Sarama configuration
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0 // Adjust based on your Kafka version

	// Create ClusterAdmin client
	admin, err := sarama.NewClusterAdmin(brokerList, config)
	if err != nil {
		return fmt.Errorf("failed to create Kafka cluster admin: %w", err)
	}
	defer admin.Close()

	// Fetch topic configurations
	topicConfigs := GetTopicConfigs()

	// Loop through each topic and create it
	for _, topicConfig := range topicConfigs {
		// Convert retention time to milliseconds
		retentionMs := fmt.Sprintf("%d", topicConfig.RetentionTime.Milliseconds())

		// Merge default configs with specific topic configs
		configs := topicConfig.Configs
		configs["retention.ms"] = retentionMs

		// Convert map[string]string to map[string]*string
		configEntries := make(map[string]*string)
		for key, value := range configs {
			val := value // Create a new variable to take the address
			configEntries[key] = &val
		}

		// Prepare Sarama TopicDetail
		detail := &sarama.TopicDetail{
			NumPartitions:     topicConfig.NumPartitions,
			ReplicationFactor: topicConfig.ReplicationFactor,
			ConfigEntries:     configEntries,
		}

		// Create the topic
		err := admin.CreateTopic(topicConfig.Name, detail, false)
		if err != nil {
			// Log if the topic already exists, continue with other topics
			if err == sarama.ErrTopicAlreadyExists {
				log.Printf("Topic %s already exists\n", topicConfig.Name)
				continue
			}
			return fmt.Errorf("failed to create topic %s: %w", topicConfig.Name, err)
		}

		log.Printf("Successfully created topic: %s\n", topicConfig.Name)
	}

	return nil
}
