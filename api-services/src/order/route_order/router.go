package route_order

import (
	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/pkg/kafka/saga"
	"github.com/supersida159/e-commerce/api-services/pkg/middleware"
	"github.com/supersida159/e-commerce/api-services/src/order/route_order/gin_order"
)

func Routes(r *gin.RouterGroup, appCtx app_context.AppContext, saga *saga.Orchestrator) {
	authRoute := r.Group("/Private", middleware.RequireAuth(appCtx))
	{
		// Create Order
		authRoute.POST("/createOrder", gin_order.CreateOrderHandler(appCtx))

		// Polling Endpoint
		authRoute.GET("/order-status/:orderId", gin_order.GetOrderStatusHandler(appCtx))

		// WebSocket Endpoint
		r.GET("/ws/order-status", gin.WrapH(gin_order.WebSocketOrderStatusHandler(appCtx, saga)))

		// Uncomment other routes as needed
		// authRoute.GET("/list", gin_order.ListOrders(appCtx))
		// authRoute.GET("/getOrder/:id", gin_order.Getorder(appCtx))
		// authRoute.POST("/softDeleteOrder", gin_order.SoftDeleteProductHandler(appCtx))
		// authRoute.PUT("/updateOrder/:id", gin_order.UpdateOrderHandler(appCtx))
	}
}
