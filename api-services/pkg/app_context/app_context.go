package app_context

import (
	"github.com/redis/go-redis/v9"
	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/config"
	dbs "github.com/supersida159/e-commerce/api-services/pkg/db"
	"github.com/supersida159/e-commerce/api-services/pkg/kafka/consumerlocal"
	"github.com/supersida159/e-commerce/api-services/pkg/kafka/producers"
	"github.com/supersida159/e-commerce/api-services/pkg/localredis"
	"github.com/supersida159/e-commerce/api-services/pkg/pubsub"
	"github.com/supersida159/e-commerce/api-services/pkg/uploadprovider"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type AppContext interface {
	GetMainDBConnection() *gorm.DB
	UploadProvider() uploadprovider.UploadProvider
	GetSecretKey() string
	GetPubSub() pubsub.PubSub
	GetCache() *localredis.RedisWRealStore
	GetRedisClient() *redis.Client
	GetConfig() *config.Schema
	GetValidatetor() *common.Validator
	GetProducer() *producers.OrderProducer
	GetConsumer() *consumerlocal.SagaConsumer
	GetOAuth() *oauth2.Config
}

type AppCtx struct {
	Dbs         *dbs.Database
	UpProvider  uploadprovider.UploadProvider
	Pb          pubsub.PubSub
	Cfg         *config.Schema
	Cache       *localredis.RedisWRealStore
	RedisClient *redis.Client
	Validator   *common.Validator
	Producer    *producers.OrderProducer
	Consumer    *consumerlocal.SagaConsumer
	OAuthConfig *oauth2.Config
}

func NewAppContext(dbs *dbs.Database, pb pubsub.PubSub, cache *localredis.RedisWRealStore, producer *producers.OrderProducer, consumer *consumerlocal.SagaConsumer, OAuthConfig *oauth2.Config) *AppCtx {
	return &AppCtx{
		Dbs: dbs,
		UpProvider: uploadprovider.NewS3Provider(config.GetConfig().AWSS3.Bucket,
			config.GetConfig().AWSS3.Region,
			config.GetConfig().AWSS3.AccessKeyID,
			config.GetConfig().AWSS3.SecretAccessKey,
			config.GetConfig().AWSS3.EndPoint),
		Pb:          pb,
		Cfg:         config.GetConfig(),
		Cache:       cache,
		Validator:   common.NewValidator(),
		Producer:    producer,
		Consumer:    consumer,
		OAuthConfig: OAuthConfig,
	}
}

func (ctx *AppCtx) GetMainDBConnection() *gorm.DB {
	return ctx.Dbs.GetDB()
}

func (ctx *AppCtx) GetSecretKey() string {
	return ctx.Cfg.AuthSecret
}

func (ctx *AppCtx) UploadProvider() uploadprovider.UploadProvider {
	return ctx.UpProvider
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

func (ctx *AppCtx) GetValidatetor() *common.Validator {
	return ctx.Validator
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

func (ctx *AppCtx) GetOAuth() *oauth2.Config {
	return ctx.OAuthConfig
}
