package httpServer

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/pkg/app_context"
	"github.com/supersida159/e-commerce/pkg/config"
	"github.com/supersida159/e-commerce/src/order/route_order"
	route_product "github.com/supersida159/e-commerce/src/product/route_product"
	"github.com/supersida159/e-commerce/src/upload/route_upload"
	"github.com/supersida159/e-commerce/src/users/route_user"
)

type Server struct {
	appCtx app_context.Appcontext
	engine *gin.Engine
}

func NewServer(appCtx app_context.Appcontext) *Server {
	return &Server{
		appCtx: appCtx,
		engine: gin.Default(),
	}
}

func (s *Server) Run() error {
	// appctx := app_context.NewAppContext(db, s3Provider, secrectkey, pubsublocal.NewPubSub())
	_ = s.engine.SetTrustedProxies(nil)
	if s.appCtx.GetConfig().Environment == config.ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}
	if err := s.MapRoutes(); err != nil {
		log.Fatalf("MapRoutes Error: %v", err)
	}

	if err := s.engine.Run(fmt.Sprintf(":%d", s.appCtx.GetConfig().HttpPort)); err != nil {
		log.Fatalf("Running HTTP server: %v", err)
	}
	return nil

}
func (s Server) GetEngine() *gin.Engine {

	return s.engine
}
func (s *Server) MapRoutes() error {

	v1 := s.engine.Group("/api/v1")

	route_user.Routes(v1.Group("/user"), s.appCtx)
	route_upload.Routes(v1.Group("/upload"), s.appCtx)
	route_product.Routes(v1.Group("/product"), s.appCtx)
	route_order.Routes(v1.Group("/order"), s.appCtx)
	return nil

}
