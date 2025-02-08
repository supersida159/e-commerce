package app_context

import (
	"github.com/redis/go-redis/v9"
	"github.com/supersida159/e-commerce/update-cart-service/pkg/config"
	dbs "github.com/supersida159/e-commerce/update-cart-service/pkg/db"
	"github.com/supersida159/e-commerce/update-cart-service/pkg/kafka/consumerlocal"
	"github.com/supersida159/e-commerce/update-cart-service/pkg/kafka/producers"
	"github.com/supersida159/e-commerce/update-cart-service/pkg/localredis"
	"github.com/supersida159/e-commerce/update-cart-service/pkg/pubsub"
	"gorm.io/gorm"
)

type AppContext interface {
	GetMainDBConnection() *gorm.DB
	GetSecretKey() string
	GetPubSub() pubsub.PubSub
	GetCache() *localredis.RedisWRealStore
	GetRedisClient() *redis.Client
	GetConfig() *config.Schema
	GetConsumer() *consumerlocal.SagaConsumer
	GetProducer() *producers.OrderProducer
}

type AppCtx struct {
	Dbs         *dbs.Database
	Pb          pubsub.PubSub
	Cfg         *config.Schema
	Cache       *localredis.RedisWRealStore
	RedisClient *redis.Client
	Consumer    *consumerlocal.SagaConsumer
	Producer    *producers.OrderProducer
}

func NewAppContext(dbs *dbs.Database, pb pubsub.PubSub, cache *localredis.RedisWRealStore, producer *producers.OrderProducer, consumer *consumerlocal.SagaConsumer) *AppCtx {
	return &AppCtx{
		Dbs:         dbs,
		Pb:          pb,
		Cfg:         config.GetConfig(),
		Cache:       cache,
		RedisClient: cache.GetClient(),
		Consumer:    consumer,
		Producer:    producer,
	}
}
func (ctx *AppCtx) GetMainDBConnection() *gorm.DB {
	return ctx.Dbs.GetDB()
}

func (ctx *AppCtx) GetSecretKey() string {
	return ctx.Cfg.AuthSecret
}

func (ctx *AppCtx) GetPubSub() pubsub.PubSub {
	return ctx.Pb
}

func (ctx *AppCtx) GetCache() *localredis.RedisWRealStore {
	return ctx.Cache
}

func (ctx *AppCtx) GetConfig() *config.Schema {
	return ctx.Cfg
}

func (ctx *AppCtx) GetProducer() *producers.OrderProducer {
	return ctx.Producer
}

func (ctx *AppCtx) GetRedisClient() *redis.Client {
	return ctx.RedisClient
}

func (ctx *AppCtx) GetConsumer() *consumerlocal.SagaConsumer {
	return ctx.Consumer
}
