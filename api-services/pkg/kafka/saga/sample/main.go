package main

import (
	"context"
	"log"
	"time"

	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/pkg/config" // Make sure to import the config package
	dbs "github.com/supersida159/e-commerce/api-services/pkg/db"
	"github.com/supersida159/e-commerce/api-services/pkg/kafka/consumerlocal"
	kafkaconfig "github.com/supersida159/e-commerce/api-services/pkg/kafka/kafka_config"
	"github.com/supersida159/e-commerce/api-services/pkg/kafka/producers"
	"github.com/supersida159/e-commerce/api-services/pkg/kafka/saga"
	"github.com/supersida159/e-commerce/api-services/pkg/localredis"
	"github.com/supersida159/e-commerce/api-services/pkg/pubsub/pubsublocal"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
)

func main() {
	cfg := config.LoadConfig()

	// Setup application context
	brokers := cfg.Kafka.Broker // Using brokers from the config
	newKafkaConfig := kafkaconfig.NewKafkaConfig(brokers)
	// Initialize Kafka producer
	producer, err := producers.NewOrderProducer(newKafkaConfig) // Pass config.Schema to producer
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Initialize Redis cache
	redisCache := localredis.NewRedis(localredis.Config{
		Address:  "localhost:6379",
		Password: "",
		Database: 0,
	}, nil)
	if err != nil {
		log.Fatalf("Failed to create Redis store: %v", err)
	}

	// Initialize PubSub (placeholder if needed)
	pubSub := pubsublocal.NewPubSub()
	dbInstance, err := dbs.NewDatabase(cfg.DatabaseURI)

	// Initialize database connection (placeholder)

	// Initialize Kafka consumer using config.Schema
	consumer, err := consumerlocal.NewSagaConsumer(newKafkaConfig, brokers, kafkaconfig.OrchestratorService) // Pass config.Schema to consumer
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	// Create application context
	appCtx := app_context.NewAppContext(dbInstance, pubSub, redisCache, producer, consumer)
	// Initialize Orchestrator
	orchestrator := saga.NewOrchestrator(appCtx)

	// Simulate a saga order event
	orderEvent := entities_orders.NewOrderEvent(&entities_orders.Order{}, entities_orders.EventSagaStarted)
	orderEvent.SagaID = "test-saga-id"
	// Start the saga process
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	go func() {
		err := consumer.Start(ctx)
		if err != nil {
			log.Printf("Consumer error: %v", err)
		}
	}()
	// Run consumer in the background
	time.Sleep(1 * time.Second)

	go func() {
		updateRollbackChannel, _ := consumer.GetEventChannel("", consumerlocal.RollbackUpdateChannel)
		m := <-updateRollbackChannel
		log.Printf("Received Update rollback event: %v", m)
	}()
	go func() {
		orchestrator.StartRollbackSingleListener(ctx)
	}()

	// err = orchestrator.StartSaga(ctx, &orderEvent)
	// if err != nil {
	// 	log.Fatalf("Failed to start saga: %v", err)
	// }

	log.Println("Saga process started successfully.")

	// Wait to observe the saga processing
	time.Sleep(30 * time.Minute)

	log.Println("Test complete.")
}
