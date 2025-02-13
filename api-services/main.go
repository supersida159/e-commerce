package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/pkg/config"
	dbs "github.com/supersida159/e-commerce/api-services/pkg/db"
	"github.com/supersida159/e-commerce/api-services/pkg/goroutineinmain"
	"github.com/supersida159/e-commerce/api-services/pkg/kafka/consumerlocal"
	kafkaconfig "github.com/supersida159/e-commerce/api-services/pkg/kafka/kafka_config"
	"github.com/supersida159/e-commerce/api-services/pkg/kafka/producers"
	"github.com/supersida159/e-commerce/api-services/pkg/kafka/saga"
	"github.com/supersida159/e-commerce/api-services/pkg/localredis"
	"github.com/supersida159/e-commerce/api-services/pkg/pubsub/pubsublocal"
	"github.com/supersida159/e-commerce/api-services/pkg/skio"
	"github.com/supersida159/e-commerce/api-services/pkg/subscriber"
	entities_carts "github.com/supersida159/e-commerce/api-services/src/cart/entities_cart"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
	"github.com/supersida159/e-commerce/api-services/src/product/entities_product"
	httpServer "github.com/supersida159/e-commerce/api-services/src/server"
	"github.com/supersida159/e-commerce/api-services/src/users/entities_user"
	"github.com/supersida159/e-commerce/api-services/src/users/repository_user"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	cfg := config.LoadConfig()

	// logger.Initialize(cfg.Environment)
	fmt.Println("url db222:", cfg.DatabaseURI)

	db, err := dbs.NewDatabase(cfg.DatabaseURI)
	fmt.Println("url db:", cfg.DatabaseURI)
	if err != nil {
		logrus.Fatal("Cannot connect to database", err)
	}
	err = db.AutoMigrate(
		&entities_user.User{},
		&entities_product.Product{},
		&entities_orders.Order{},
		&entities_carts.Cart{},
		&common.Image{},
		&entities_product.CartItem{},
		&entities_user.Address{},
	)

	if err != nil {
		logrus.Fatal(" Cannot connect to database to AutoMigrate", err)
	}

	cache := localredis.NewRedis(localredis.Config{
		Address:  cfg.RedisURI,
		Password: cfg.RedisPassword,
		Database: cfg.RedisDB,
	},
		repository_user.NewSQLStore(db.GetDB()),
	)
	connectRedis := cache.IsConnected()
	fmt.Println("connect redis:", connectRedis)
	localpubsub := pubsublocal.NewPubSub()

	brokers := cfg.Kafka.Broker
	newKafkaConfig := kafkaconfig.NewKafkaConfig(brokers)
	producer, err := producers.NewOrderProducer(newKafkaConfig) // Pass config.Schema to producer
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Initialize Kafka consumer using config.Schema
	consumer, err := consumerlocal.NewSagaConsumer(newKafkaConfig, brokers, kafkaconfig.OrchestratorService) // Pass config.Schema to consumer
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	conf := &oauth2.Config{
		ClientID:     cfg.OAuth.ClientID,
		ClientSecret: cfg.OAuth.ClientSecret,
		RedirectURL:  "http://localhost:8090/api/v1/auth/callback",
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	appctx := app_context.NewAppContext(db, localpubsub, cache, producer, consumer, conf)

	err = goroutineinmain.RunExpireOrder(appctx)
	if err != nil {
		logrus.Fatal(" Cannot connect to database to AutoMigrate", err)
	}

	rtengine := skio.NewEngine()
	newOrchestrator := saga.NewOrchestrator(appctx)
	if err := subscriber.NewEngine(appctx, newOrchestrator).Start(); err != nil {
		log.Fatalln(err)
	}
	// Start consumer in main
	ctx := context.Background()
	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Fatalf("Consumer failed: %v", err)
		}
	}()

	httpSvr := httpServer.NewServer(appctx, newOrchestrator)

	httpSvr.GetEngine().Use(CORSMiddleware())
	if err = httpSvr.Run(rtengine); err != nil {
		logrus.Fatal(" Cannot runHttp server", err)
	}

}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		fmt.Print(c.Request.Method)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
