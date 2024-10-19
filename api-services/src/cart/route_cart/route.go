package route_cart

import (
	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/pkg/middleware"
	gin_carts "github.com/supersida159/e-commerce/api-services/src/cart/route_cart/gin_cart"
)

func Routes(r *gin.RouterGroup, appCtx app_context.Appcontext) {
	ginCarts, _ := gin_carts.NewCartController(appCtx)
	// r.POST("/register", gin_user.Register(appCtx))
	authRoute := r.Group("/Private", middleware.RequireAuth(appCtx))
	{
		authRoute.POST("/createCart", ginCarts.CreateCart)
		authRoute.DELETE("/deleteCart", ginCarts.DeleteCart)
		authRoute.GET("/getCart", ginCarts.GetCart)
		// authRoute.POST("/softDeleteProduct", gin_product.SoftDeleteProductHandler(appCtx))
		// authRoute.POST("/updateProduct", gin_product.UpdateProductHandler(appCtx))
		authRoute.PUT("/updateCart", ginCarts.UpdateCart)
		// authRoute.POST("/adminUpdate", gin_user.AddUpdateAddmin(appCtx))
		// authRoute.POST("/register", userHandler.Register)
		// authRoute.POST("/login", userHandler.Login)
		// authRoute.POST("/refresh", refreshAuthMiddleware, userHandler.RefreshToken)
		// authRoute.GET("/me", authMiddleware, userHandler.GetMe)
		// authRoute.PUT("/change-password", authMiddleware, userHandler.ChangePassword)
	}
}
