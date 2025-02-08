package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/supersida159/e-commerce/update-cart-service/common"
	"github.com/supersida159/e-commerce/update-cart-service/pkg/app_context"
	"github.com/supersida159/e-commerce/update-cart-service/pkg/config"
	dbs "github.com/supersida159/e-commerce/update-cart-service/pkg/db"
	"github.com/supersida159/e-commerce/update-cart-service/pkg/kafka/consumerlocal"
	kafkaconfig "github.com/supersida159/e-commerce/update-cart-service/pkg/kafka/kafka_config"
	"github.com/supersida159/e-commerce/update-cart-service/pkg/kafka/producers"
	"github.com/supersida159/e-commerce/update-cart-service/pkg/localredis"
	"github.com/supersida159/e-commerce/update-cart-service/pkg/pubsub/pubsublocal"
	"github.com/supersida159/e-commerce/update-cart-service/pkg/subscriber"
	entities "github.com/supersida159/e-commerce/update-cart-service/src/model"
	"github.com/supersida159/e-commerce/update-cart-service/src/repository"
	"github.com/supersida159/e-commerce/update-cart-service/src/usecase"
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
	cartStore := repository.NewSQLStore(db.GetDB())

	// Initialize Redis cache
	cache := localredis.NewRedis(localredis.Config{
		Address:  cfg.RedisURI,
		Password: cfg.RedisPassword,
		Database: cfg.RedisDB,
	}, cartStore)

	// Initialize local pubsub
	localPubSub := pubsublocal.NewPubSub()

	// Initialize Kafka consumer
	consumer, err := consumerlocal.NewSagaConsumer(cfg, cfg.Kafka.Brokers, kafkaconfig.CartService)
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
	cartUseCase := usecase.NewCartUsecase(cartStore)

	// Initialize subscriber
	sub := subscriber.NewSubscriber(appCtx)

	// Register handlers
	sub.RegisterHandler(consumerlocal.CreateOrderChannel, func(ctx context.Context, data *entities.OrderEvent) *common.AppError {
		result := cartUseCase.SoftDeleteCart(ctx, data.CartID)
		data.CurrentService = entities.CartService
		if result != nil {

			data.ServiceStatus.LastUpdated = time.Now()
			data.ServiceStatus.ServiceStates[entities.CartService] = entities.ServiceState{
				Status:    entities.ServiceFailed,
				UpdatedAt: time.Now(),
				Error:     result.RootErr.Error(),
			}
		} else {
			data.ServiceStatus.LastUpdated = time.Now()
			data.ServiceStatus.ServiceStates[entities.CartService] = entities.ServiceState{
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

		result := cartUseCase.RecoveryCart(ctx, data.CartID)
		data.CurrentService = entities.CartService

		if result != nil {
			data.ServiceStatus.LastUpdated = time.Now()
			data.ServiceStatus.ServiceStates[entities.CartService] = entities.ServiceState{
				Status:    entities.ServiceRollbackFailed,
				UpdatedAt: time.Now(),
				Error:     result.RootErr.Error(),
			}
		} else {
			data.ServiceStatus.LastUpdated = time.Now()
			data.ServiceStatus.ServiceStates[entities.CartService] = entities.ServiceState{
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
