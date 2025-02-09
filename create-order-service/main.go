package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/supersida159/e-commerce/create-order/common"
	"github.com/supersida159/e-commerce/create-order/pkg/app_context"
	"github.com/supersida159/e-commerce/create-order/pkg/config"
	dbs "github.com/supersida159/e-commerce/create-order/pkg/db"
	"github.com/supersida159/e-commerce/create-order/pkg/kafka/consumerlocal"
	kafkaconfig "github.com/supersida159/e-commerce/create-order/pkg/kafka/kafka_config"
	"github.com/supersida159/e-commerce/create-order/pkg/kafka/producers"
	"github.com/supersida159/e-commerce/create-order/pkg/localredis"
	"github.com/supersida159/e-commerce/create-order/pkg/pubsub/pubsublocal"
	"github.com/supersida159/e-commerce/create-order/pkg/subscriber"
	entities "github.com/supersida159/e-commerce/create-order/src/model"
	repository_orders "github.com/supersida159/e-commerce/create-order/src/repository"
	usecase_orders "github.com/supersida159/e-commerce/create-order/src/usecase"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := dbs.NewDatabase(cfg.DatabaseURI)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize repository
	orderStore := repository_orders.NewSQLStore(db.GetDB())

	// Initialize Redis cache
	cache := localredis.NewRedis(localredis.Config{
		Address:  cfg.RedisURI,
		Password: cfg.RedisPassword,
		Database: cfg.RedisDB,
	}, orderStore)

	// Initialize local pubsub
	localPubSub := pubsublocal.NewPubSub()

	// Initialize Kafka consumer
	consumer, err := consumerlocal.NewSagaConsumer(cfg, cfg.Kafka.Brokers, kafkaconfig.OrderService)
	if err != nil {
		log.Fatal("Failed to initialize Kafka consumer:", err)
	}
	defer consumer.Stop()

	// Initialize Kafka producer
	producer, err := producers.NewOrderProducer(cfg)
	if err != nil {
		log.Fatal("Failed to initialize Kafka producer:", err)
	}
	defer producer.Close()

	// Initialize application context
	appCtx := app_context.NewAppContext(db, localPubSub, cache, producer, consumer)

	// Initialize order usecase
	orderUseCase := usecase_orders.NewOrderUsecase(orderStore)

	// Initialize subscriber
	sub := subscriber.NewSubscriber(appCtx)

	// Register handlers
	sub.RegisterHandler(consumerlocal.CreateOrderChannel, func(ctx context.Context, data *entities.OrderEvent) *common.AppError {
		orderId, result := orderUseCase.CreateOrder(ctx, &data.Order)
		data.CurrentService = entities.OrderService
		if result != nil {

			data.ServiceStatus.LastUpdated = time.Now()
			data.ServiceStatus.ServiceStates[entities.OrderService] = entities.ServiceState{
				Status:    entities.ServiceFailed,
				UpdatedAt: time.Now(),
				Error:     result.RootErr.Error(),
			}
		} else {
			data.Order.ID = *orderId
			data.ServiceStatus.LastUpdated = time.Now()
			data.ServiceStatus.ServiceStates[entities.OrderService] = entities.ServiceState{
				Status:    entities.ServiceSuccess,
				UpdatedAt: time.Now(),
				Error:     "",
			}
		}
		err := sub.GetAppContext().GetProducer().SendStatusUpdate(*data)
		if err != nil {
			log.Fatal("Failed to send status update:", err)
		}
		return nil
	})

	sub.RegisterHandler(consumerlocal.RollbackChannel, func(ctx context.Context, data *entities.OrderEvent) *common.AppError {

		result := orderUseCase.SoftDeleteOrder(ctx, data.Order.ID)
		data.CurrentService = entities.OrderService

		if result != nil {
			data.ServiceStatus.LastUpdated = time.Now()
			data.ServiceStatus.ServiceStates[entities.OrderService] = entities.ServiceState{
				Status:    entities.ServiceRollbackFailed,
				UpdatedAt: time.Now(),
				Error:     result.RootErr.Error(),
			}
		} else {
			data.ServiceStatus.LastUpdated = time.Now()
			data.ServiceStatus.ServiceStates[entities.OrderService] = entities.ServiceState{
				Status:    entities.ServiceRollbackSuccess,
				UpdatedAt: time.Now(),
				Error:     "",
			}
		}

		err := sub.GetAppContext().GetProducer().SendUpdateRollback(*data)

		if err != nil {
			log.Fatal("Failed to send status update:", err)
		}

		return nil

	})

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start services
	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Fatal("Failed to start consumer:", err)
		}
	}()

	go func() {
		if err := sub.Start(ctx); err != nil {
			log.Fatal("Failed to start subscriber:", err)
		}
	}()

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down gracefully...")

	// Clean shutdown
	sub.Stop()
	consumer.Stop()
}
