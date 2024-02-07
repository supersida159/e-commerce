package app_context

import (
	"github.com/supersida159/e-commerce/pkg/config"
	dbs "github.com/supersida159/e-commerce/pkg/db"
	"github.com/supersida159/e-commerce/pkg/pubsub"
	"github.com/supersida159/e-commerce/pkg/redis"
	"github.com/supersida159/e-commerce/pkg/uploadprovider"
	"gorm.io/gorm"
)

type Appcontext interface {
	GetMainDBConnection() *gorm.DB
	UploadProvider() uploadprovider.UploadProvider
	GetSecretKey() string
	GetPubSub() pubsub.PubSub
	GetCache() *redis.RedisWRealStore
	GetConfig() *config.Schema
}

type AppCtx struct {
	Dbs        *dbs.Database
	UpProvider uploadprovider.UploadProvider
	Pb         pubsub.PubSub
	Cfg        *config.Schema
	Cache      *redis.RedisWRealStore
}

func NewAppContext(dbs *dbs.Database, pb pubsub.PubSub, cache *redis.RedisWRealStore) *AppCtx {
	return &AppCtx{
		Dbs:        dbs,
		UpProvider: uploadprovider.NewS3Provider(config.GetConfig().S3BucketName, config.GetConfig().S3Region, config.GetConfig().S3APIKey, config.GetConfig().S3SecretKey, config.GetConfig().S3Domain),
		Pb:         pb,
		Cfg:        config.GetConfig(),
		Cache:      cache,
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

func (ctx *AppCtx) GetCache() *redis.RedisWRealStore {
	return ctx.Cache
}

func (ctx *AppCtx) GetConfig() *config.Schema {
	return ctx.Cfg
}
