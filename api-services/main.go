package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/pkg/config"
	dbs "github.com/supersida159/e-commerce/api-services/pkg/db"
	"github.com/supersida159/e-commerce/api-services/pkg/goroutineinmain"
	"github.com/supersida159/e-commerce/api-services/pkg/pubsub/pubsublocal"
	"github.com/supersida159/e-commerce/api-services/pkg/redis"
	"github.com/supersida159/e-commerce/api-services/pkg/skio"
	entities_carts "github.com/supersida159/e-commerce/api-services/src/cart/entities_cart"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
	"github.com/supersida159/e-commerce/api-services/src/product/entities_product"
	httpServer "github.com/supersida159/e-commerce/api-services/src/server"
	subscriber "github.com/supersida159/e-commerce/api-services/src/subcriber"
	"github.com/supersida159/e-commerce/api-services/src/users/entities_user"
	"github.com/supersida159/e-commerce/api-services/src/users/repository_user"
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

	cache := redis.NewRedis(redis.Config{
		Address:  cfg.RedisURI,
		Password: cfg.RedisPassword,
		Database: cfg.RedisDB,
	},
		repository_user.NewSQLStore(db.GetDB()),
	)
	connectRedis := cache.IsConnected()
	fmt.Println("connect redis:", connectRedis)
	localpubsub := pubsublocal.NewPubSub()
	appctx := app_context.NewAppContext(db, localpubsub, cache)

	err = goroutineinmain.RunExpireOrder(appctx)
	if err != nil {
		logrus.Fatal(" Cannot connect to database to AutoMigrate", err)
	}

	rtengine := skio.NewEngine()

	if err := subscriber.NewEngine(appctx, rtengine).Start(); err != nil {
		log.Fatalln(err)
	}

	httpSvr := httpServer.NewServer(appctx)

	httpSvr.GetEngine().Use(CORSMiddleware())
	if err = httpSvr.Run(rtengine); err != nil {
		logrus.Fatal(" Cannot runHttp server", err)
	}

	// run below code with Grpc

	// go func() {
	// 	httpSvr := httpServer.NewServer(appctx)
	// 	if err = httpSvr.Run(); err != nil {
	// 		logrus.Fatal(err)
	// 	}
	// }()
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
