package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/supersida159/e-commerce/pkg/app_context"
	"github.com/supersida159/e-commerce/pkg/config"
	dbs "github.com/supersida159/e-commerce/pkg/db"
	"github.com/supersida159/e-commerce/pkg/pubsub/pubsublocal"
	"github.com/supersida159/e-commerce/pkg/redis"
	"github.com/supersida159/e-commerce/src/product/entities_product"
	httpServer "github.com/supersida159/e-commerce/src/server"
	"github.com/supersida159/e-commerce/src/users/entities"
	"github.com/supersida159/e-commerce/src/users/repository_user"
)

func main() {
	cfg := config.LoadConfig()
	// logger.Initialize(cfg.Environment)

	db, err := dbs.NewDatabase(cfg.DatabaseURI)
	fmt.Println("url db:", cfg.DatabaseURI)
	if err != nil {
		logrus.Fatal("Cannot connect to database", err)
	}
	err = db.AutoMigrate(&entities.User{}, &entities_product.Product{})
	if err != nil {
		logrus.Fatal(err)
	}
	cache := redis.NewRedis(redis.Config{
		Address:  cfg.RedisURI,
		Password: cfg.RedisPassword,
		Database: cfg.RedisDB,
	},
		repository_user.NewSQLStore(db.GetDB()),
	)
	connectRedis := cache.IsConnected()
	fmt.Println("connect redis:", connectRedis)
	appctx := app_context.NewAppContext(db, pubsublocal.NewPubSub(), cache)

	httpSvr := httpServer.NewServer(appctx)
	if err = httpSvr.Run(); err != nil {
		logrus.Fatal(err)
	}

	// run below code with Grpc

	// go func() {
	// 	httpSvr := httpServer.NewServer(appctx)
	// 	if err = httpSvr.Run(); err != nil {
	// 		logrus.Fatal(err)
	// 	}
	// }()
}
