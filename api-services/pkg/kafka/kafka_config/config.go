package kafkaconfig

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"golang.org/x/sync/errgroup"
)

// ServiceID represents different services in the ecosystem
type ServiceID string

const (
	OrchestratorService ServiceID = "orchestrator"
	OrderService        ServiceID = "order"
	InventoryService    ServiceID = "inventory"
	CartService         ServiceID = "cart"
)

var KafkaTopics = struct {
	CreateOrder       string
	UpdateOrder       string
	UpdateRollback    string
	RollbackOrder     string
	RollbackCart      string
	RollbackInventory string
}{
	CreateOrder:       "saga.create.order",
	UpdateOrder:       "saga.update.order",
	UpdateRollback:    "saga.update.rollback",
	RollbackOrder:     "saga.rollback.order",
	RollbackCart:      "saga.rollback.cart",
	RollbackInventory: "saga.rollback.inventory",
}

// TopicConfig defines complete Kafka topic configuration
type TopicConfig struct {
	Name              string
	Partitions        int32
	ReplicationFactor int16
	Configs           []TopicProperty
}

type TopicProperty struct {
	Key   string
	Value string
}

// KafkaConfig wraps all Kafka configuration parameters
type KafkaConfig struct {
	Brokers    []string
	Producer   *sarama.Config
	Consumer   *sarama.Config
	admin      sarama.ClusterAdmin
	adminOnce  sync.Once
	adminMutex sync.Mutex
}

var (
	defaultRetention = 7 * 24 * time.Hour
	topicRegistry    = map[string]TopicConfig{
		KafkaTopics.CreateOrder: {
			Partitions:        6,
			ReplicationFactor: 3,
			Configs: []TopicProperty{
				{Key: "cleanup.policy", Value: "compact"},
				{Key: "retention.ms", Value: fmt.Sprintf("%d", defaultRetention.Milliseconds())},
			},
		},
		KafkaTopics.UpdateOrder: {
			Partitions:        1,
			ReplicationFactor: 3,
			Configs: []TopicProperty{
				{Key: "cleanup.policy", Value: "delete"},
				{Key: "retention.ms", Value: fmt.Sprintf("%d", defaultRetention.Milliseconds())},
			},
		},
		KafkaTopics.RollbackOrder:     createRollbackTopicConfig(),
		KafkaTopics.RollbackInventory: createRollbackTopicConfig(),
		KafkaTopics.RollbackCart:      createRollbackTopicConfig(),
		KafkaTopics.UpdateRollback:    createRollbackTopicConfig(),
	}
)

// NewKafkaConfig creates a new validated Kafka configuration
func NewKafkaConfig(brokers []string) *KafkaConfig {
	return &KafkaConfig{
		Brokers:  brokers,
		Producer: createProducerConfig(),
		Consumer: createConsumerConfig(""),
	}
}

// In your kafkaconfig package's producer configuration
func createProducerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V3_4_0_0

	// Required for SyncProducer
	config.Producer.Return.Successes = true // 👈 Add this line

	// Idempotent producer settings
	config.Producer.Idempotent = true
	config.Net.MaxOpenRequests = 1

	// Exactly-once semantics requirements
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Retry.Backoff = 1 * time.Second

	// Optional optimizations
	config.Producer.Compression = sarama.CompressionSnappy
	return config
}

// SaramaConfig returns the appropriate Sarama configuration
func (kc *KafkaConfig) SaramaConfig(isProducer bool) *sarama.Config {
	if isProducer {
		return kc.Producer
	}
	return kc.Consumer
}

// InitializeTopology creates all required topics
func (kc *KafkaConfig) InitializeTopology(ctx context.Context) error {
	admin, err := kc.getClusterAdmin()
	if err != nil {
		return err
	}
	defer kc.closeAdmin()

	existingTopics, err := admin.ListTopics()
	if err != nil {
		return fmt.Errorf("failed to list existing topics: %w", err)
	}

	var (
		g, _     = errgroup.WithContext(ctx)
		mu       sync.Mutex
		errSlice []error
	)

	for name, cfg := range topicRegistry {
		if _, exists := existingTopics[name]; exists {
			log.Printf("Topic %s already exists", name)
			continue
		}

		name, cfg := name, cfg
		g.Go(func() error {
			detail := sarama.TopicDetail{
				NumPartitions:     cfg.Partitions,
				ReplicationFactor: cfg.ReplicationFactor,
				ConfigEntries:     make(map[string]*string),
			}

			for _, prop := range cfg.Configs {
				val := prop.Value
				detail.ConfigEntries[prop.Key] = &val
			}

			if err := admin.CreateTopic(name, &detail, false); err != nil {
				mu.Lock()
				errSlice = append(errSlice, fmt.Errorf("failed to create topic %s: %w", name, err))
				mu.Unlock()
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("topic creation failed: %w", err)
	}

	if len(errSlice) > 0 {
		return fmt.Errorf("topic creation errors: %v", errSlice)
	}

	return nil
}

// GetServiceTopics returns topics for a specific service
func GetServiceTopics(service ServiceID) []string {
	switch service {
	case OrchestratorService:
		return []string{"saga.update.order", "saga.update.rollback"}
	case OrderService:
		return []string{"saga.create.order", "saga.rollback.order"}
	case InventoryService:
		return []string{"saga.create.order", "saga.rollback.inventory"}
	case CartService:
		return []string{"saga.create.order", "saga.rollback.cart"}
	default:
		return nil
	}
}

// ConsumerGroupName generates consumer group ID for a service
func ConsumerGroupName(service ServiceID) string {
	return fmt.Sprintf("%s-consumer-group", service)
}

func createRollbackTopicConfig() TopicConfig {
	return TopicConfig{
		Partitions:        3,
		ReplicationFactor: 3,
		Configs: []TopicProperty{
			{Key: "cleanup.policy", Value: "compact"},
			{Key: "retention.ms", Value: fmt.Sprintf("%d", (30 * 24 * time.Hour).Milliseconds())},
			{Key: "min.compaction.lag.ms", Value: "3600000"},
		},
	}
}

func createConsumerConfig(groupID string) *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V3_4_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.BalanceStrategySticky,
	}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 5 * time.Second
	config.ClientID = groupID
	return config
}

func (kc *KafkaConfig) getClusterAdmin() (sarama.ClusterAdmin, error) {
	var err error
	kc.adminOnce.Do(func() {
		kc.adminMutex.Lock()
		defer kc.adminMutex.Unlock()
		kc.admin, err = sarama.NewClusterAdmin(kc.Brokers, kc.Producer)
	})
	return kc.admin, err
}

func (kc *KafkaConfig) closeAdmin() {
	kc.adminMutex.Lock()
	defer kc.adminMutex.Unlock()
	if kc.admin != nil {
		kc.admin.Close()
		kc.admin = nil
	}
	kc.adminOnce = sync.Once{}
}

// Validate checks configuration validity
func (kc *KafkaConfig) Validate() error {
	if len(kc.Brokers) == 0 {
		return fmt.Errorf("at least one broker required")
	}
	if kc.Producer == nil || kc.Consumer == nil {
		return fmt.Errorf("missing producer/consumer configuration")
	}
	return nil
}
