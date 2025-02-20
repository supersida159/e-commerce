// route_auth/routes.go
package route_auth

import (
	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/src/OAUTH2.0/route/gin_auth"
	"github.com/supersida159/e-commerce/api-services/src/order/route_order/gin_order"
)

func Routes(r *gin.RouterGroup, appCtx app_context.AppContext) {
	// Public routes
	r.POST("/google", gin_auth.HandleGoogleLogin(appCtx))

	// Protected routes
	protected := r.Group("")
	protected.Use(gin_auth.AuthMiddleware(appCtx))
	{
		r.GET("/order-status/:orderId", gin_order.GetOrderStatusHandler(appCtx))
		// Add other protected routes here
	}
}
