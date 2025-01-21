package consumerlocal

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"sync"

// 	"github.com/IBM/sarama"
// 	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
// 	kafkaconfig "github.com/supersida159/e-commerce/api-services/pkg/kafka/kafka_config"
// 	"github.com/supersida159/e-commerce/api-services/pkg/kafka/producers"
// 	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
// )

// type SagaConsumer struct {
// 	consumer   sarama.Consumer
// 	Producer   producers.OrderProducer
// 	serviceID  kafkaconfig.ServiceID
// 	topics     []string
// 	ready      chan bool
// 	stopChan   chan struct{}
// 	wg         sync.WaitGroup
// 	stateMutex sync.RWMutex
// 	appCtx     app_context.Appcontext
// }

// func NewSagaConsumer(
// 	config *sarama.Config,
// 	brokers []string,
// 	serviceID kafkaconfig.ServiceID,
// 	producer producers.OrderProducer,
// 	appCtx app_context.Appcontext,
// ) (*SagaConsumer, error) {
// 	consumer, err := sarama.NewConsumer(brokers, config)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create consumer: %w", err)
// 	}

// 	topicMapping := kafkaconfig.GetServiceTopicMapping()
// 	topics := topicMapping[serviceID]

// 	return &SagaConsumer{
// 		consumer:  consumer,
// 		Producer:  producer,
// 		serviceID: serviceID,
// 		topics:    topics,
// 		ready:     make(chan bool),
// 		stopChan:  make(chan struct{}),
// 		appCtx:    appCtx,
// 	}, nil
// }

// func (c *SagaConsumer) Start(ctx context.Context) error {
// 	log.Printf("Starting consumer for service %s, listening to topics: %v", c.serviceID, c.topics)

// 	for _, topic := range c.topics {
// 		c.wg.Add(1)
// 		go func(topic string) {
// 			defer c.wg.Done()
// 			c.consumeTopic(ctx, topic)
// 		}(topic)
// 	}

// 	return nil
// }

// func (c *SagaConsumer) Stop() error {
// 	close(c.stopChan)
// 	c.wg.Wait()
// 	return c.consumer.Close()
// }

// func (c *SagaConsumer) consumeTopic(ctx context.Context, topic string) {
// 	partitions, err := c.consumer.Partitions(topic)
// 	if err != nil {
// 		log.Printf("failed to get partitions for topic %s: %v", topic, err)
// 		return
// 	}

// 	for _, partition := range partitions {
// 		pc, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
// 		if err != nil {
// 			log.Printf("failed to start consumer for topic %s partition %d: %v", topic, partition, err)
// 			continue
// 		}

// 		defer pc.Close()

// 		for {
// 			select {
// 			case msg := <-pc.Messages():
// 				if err := c.processMessage(ctx, msg); err != nil {
// 					log.Printf("Error processing message: %v", err)
// 				}
// 			case err := <-pc.Errors():
// 				log.Printf("Error from consumer: %v", err)
// 			case <-c.stopChan:
// 				return
// 			case <-ctx.Done():
// 				return
// 			}
// 		}
// 	}
// }

// func (c *SagaConsumer) processMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
// 	var event entities_orders.OrderEvent
// 	if err := json.Unmarshal(msg.Value, &event); err != nil {
// 		return fmt.Errorf("failed to unmarshal message: %w", err)
// 	}

// 	log.Printf("[%s] Received event type %s for saga %s", c.serviceID, event.EventType, event.SagaID)

// 	// Instead of handling logic here, emit an event that the orchestrator can handle
// 	opts := &producers.SendMessageOptions{
// 		Topic: msg.Topic,
// 	}

// 	// // Send the event to the orchestrator with the current service's status
// 	// event.ServiceStatus = entities_orders.ServiceStatus{
// 	// 	ServiceID: c.serviceID,
// 	// 	Status:    entities_orders.ServiceReceived,
// 	// }

// 	_, err := c.Producer.SendParallel(event, *opts)
// 	if err != nil {
// 		return fmt.Errorf("failed to send event to orchestrator: %w", err)
// 	}

// 	return nil
// }

// func NewDefaultSagaConsumer(
// 	brokers []string,
// 	serviceID kafkaconfig.ServiceID,
// 	producer producers.OrderProducer,
// 	appCtx app_context.Appcontext,
// ) (*SagaConsumer, error) {
// 	config := sarama.NewConfig()
// 	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
// 	config.Consumer.Offsets.Initial = sarama.OffsetNewest
// 	config.Producer.Return.Successes = true

// 	return NewSagaConsumer(config, brokers, serviceID, producer, appCtx)
// }
