package httpServer

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/pkg/config"
	"github.com/supersida159/e-commerce/api-services/pkg/kafka/saga"
	"github.com/supersida159/e-commerce/api-services/pkg/skio"
	route_auth "github.com/supersida159/e-commerce/api-services/src/OAUTH2.0/route"
	"github.com/supersida159/e-commerce/api-services/src/cart/route_cart"
	"github.com/supersida159/e-commerce/api-services/src/order/route_order"
	route_product "github.com/supersida159/e-commerce/api-services/src/product/route_product"
	"github.com/supersida159/e-commerce/api-services/src/upload/route_upload"

	// "github.com/supersida159/e-commerce/api-services/src/upload/route_upload"
	"github.com/supersida159/e-commerce/api-services/src/users/route_user"
)

type Server struct {
	appCtx app_context.AppContext
	engine *gin.Engine
	saga   *saga.Orchestrator
}

func (s Server) GetAppContext() app_context.AppContext {
	return s.appCtx
}
func NewServer(appCtx app_context.AppContext, saga *saga.Orchestrator) *Server {
	return &Server{
		appCtx: appCtx,
		engine: gin.Default(),
		saga:   saga,
	}
}

func (s *Server) Run(rte *skio.RtEngine) error {
	// appctx := app_context.NewAppContext(db, s3Provider, secrectkey, pubsublocal.NewPubSub())
	_ = s.engine.SetTrustedProxies(nil)
	if s.appCtx.GetConfig().Environment == config.ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}
	if err := s.MapRoutes(); err != nil {
		log.Fatalf("MapRoutes Error: %v", err)
	}

	rte.Run(s.appCtx, s.engine)

	fmt.Println("port:", s.appCtx.GetConfig().HttpPort)

	if err := s.engine.Run(fmt.Sprintf(":%d", s.appCtx.GetConfig().HttpPort)); err != nil {
		log.Fatalf("Running HTTP server: %v", err)
	}
	return nil

}
func (s Server) GetEngine() *gin.Engine {

	return s.engine
}
func (s *Server) MapRoutes() error {

	s.engine.StaticFile("/demo/", "./demo.html")

	// s.engine.GET("/auth/callback", gin_auth.OAuthCallbackHandler(s.appCtx))

	v1 := s.engine.Group("/api/v1")

	route_user.Routes(v1.Group("/user"), s.appCtx)
	// route_upload.Routes(v1.Group("/upload"), s.appCtx)
	route_product.Routes(v1.Group("/product"), s.appCtx)
	route_order.Routes(v1.Group("/order"), s.appCtx, s.saga)
	route_cart.Routes(v1.Group("/cart"), s.appCtx)
	route_upload.Routes(v1.Group("/upload"), s.appCtx)

	route_auth.Routes(v1.Group("/auth"), s.appCtx)

	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	return nil

}
