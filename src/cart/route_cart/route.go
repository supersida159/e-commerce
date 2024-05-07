package route_cart

import (
	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/pkg/app_context"
	"github.com/supersida159/e-commerce/pkg/middleware"
	gin_carts "github.com/supersida159/e-commerce/src/cart/route_cart/gin_cart"
)

func Routes(r *gin.RouterGroup, appCtx app_context.Appcontext) {

	// r.POST("/register", gin_user.Register(appCtx))
	authRoute := r.Group("/Private", middleware.RequireAuth(appCtx))
	{
		authRoute.POST("/createCart", gin_carts.CreateCart(appCtx))
		authRoute.DELETE("/deleteCart", gin_carts.DeleteCart(appCtx))
		authRoute.GET("/getCart", gin_carts.GetCart(appCtx))
		// authRoute.POST("/softDeleteProduct", gin_product.SoftDeleteProductHandler(appCtx))
		// authRoute.POST("/updateProduct", gin_product.UpdateProductHandler(appCtx))
		authRoute.PUT("/updateCart", gin_carts.UpdateCart(appCtx))
		// authRoute.POST("/adminUpdate", gin_user.AddUpdateAddmin(appCtx))
		// authRoute.POST("/register", userHandler.Register)
		// authRoute.POST("/login", userHandler.Login)
		// authRoute.POST("/refresh", refreshAuthMiddleware, userHandler.RefreshToken)
		// authRoute.GET("/me", authMiddleware, userHandler.GetMe)
		// authRoute.PUT("/change-password", authMiddleware, userHandler.ChangePassword)
	}
}
