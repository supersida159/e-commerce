package routeproduct

import (
	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/pkg/middleware"
	gin_product "github.com/supersida159/e-commerce/api-services/src/product/route_product/gin_product"
)

func Routes(r *gin.RouterGroup, appCtx app_context.Appcontext) {
	r.POST("/list", gin_product.ListProducts(appCtx))
	r.GET("/getProduct/:name", gin_product.GetProductHandler(appCtx))
	// r.POST("/register", gin_user.Register(appCtx))
	authRoute := r.Group("/Private", middleware.RequireAuth(appCtx))
	{
		authRoute.POST("/createProduct", gin_product.CreateProductHandler(appCtx))
		authRoute.POST("/softDeleteProduct", gin_product.SoftDeleteProductHandler(appCtx))
		authRoute.PUT("/updateProduct", gin_product.UpdateProductHandler(appCtx))
		// authRoute.PUT("/update", gin_user.UpdateUser(appCtx))
		// authRoute.POST("/adminUpdate", gin_user.AddUpdateAddmin(appCtx))
		// authRoute.POST("/register", userHandler.Register)
		// authRoute.POST("/login", userHandler.Login)
		// authRoute.POST("/refresh", refreshAuthMiddleware, userHandler.RefreshToken)
		// authRoute.GET("/me", authMiddleware, userHandler.GetMe)
		// authRoute.PUT("/change-password", authMiddleware, userHandler.ChangePassword)
	}
}
