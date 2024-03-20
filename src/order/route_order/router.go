package route_order

import (
	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/pkg/app_context"
	"github.com/supersida159/e-commerce/pkg/middleware"
	"github.com/supersida159/e-commerce/src/order/route_order/gin_order"
	gin_product "github.com/supersida159/e-commerce/src/product/route_product/gin_product"
	"github.com/supersida159/e-commerce/src/users/route_user/gin_user"
)

func Routes(r *gin.RouterGroup, appCtx app_context.Appcontext) {
	r.POST("/list", gin_order.ListOrders(appCtx))
	r.GET("/getProduct/:id", gin_user.Login(appCtx))
	// r.POST("/register", gin_user.Register(appCtx))
	authRoute := r.Group("/Private", middleware.RequireAuth(appCtx))
	{
		authRoute.POST("/createOrder", gin_order.CreateOrderHandler(appCtx))
		authRoute.POST("/softDeleteProduct", gin_product.SoftDeleteProductHandler(appCtx))
		authRoute.POST("/updateProduct", gin_product.UpdateProductHandler(appCtx))
		// authRoute.PUT("/update", gin_user.UpdateUser(appCtx))
		// authRoute.POST("/adminUpdate", gin_user.AddUpdateAddmin(appCtx))
		// authRoute.POST("/register", userHandler.Register)
		// authRoute.POST("/login", userHandler.Login)
		// authRoute.POST("/refresh", refreshAuthMiddleware, userHandler.RefreshToken)
		// authRoute.GET("/me", authMiddleware, userHandler.GetMe)
		// authRoute.PUT("/change-password", authMiddleware, userHandler.ChangePassword)
	}
}
