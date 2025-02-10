package route_auth

import (
	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/src/order/route_order/gin_order"
)

func Routes(r *gin.RouterGroup, appCtx app_context.AppContext) {

	r.GET("/login", gin_order.CreateOrderHandler(appCtx))
	r.GET("/oauth", gin_order.CreateOrderHandler(appCtx))

	r.GET("/callback", gin_order.CreateOrderHandler(appCtx))

	// Polling Endpoint
	r.GET("/order-status/:orderId", gin_order.GetOrderStatusHandler(appCtx))

}
