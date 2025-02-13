package route_auth

import (
	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/src/OAUTH2.0/route/gin_auth"
	"github.com/supersida159/e-commerce/api-services/src/order/route_order/gin_order"
)

func Routes(r *gin.RouterGroup, appCtx app_context.AppContext) {

	r.GET("/login", gin_order.CreateOrderHandler(appCtx))
	r.GET("/auth", gin_auth.OAuthHandler(appCtx))
	r.POST("/oauth", gin_auth.HandleGoogleOAuth(appCtx))

	r.GET("/callback", gin_auth.OAuthCallbackHandler(appCtx))

	// Polling Endpoint
	r.GET("/order-status/:orderId", gin_order.GetOrderStatusHandler(appCtx))

}
