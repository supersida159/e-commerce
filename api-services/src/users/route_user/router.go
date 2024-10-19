package route_user

import (
	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/pkg/middleware"
	"github.com/supersida159/e-commerce/api-services/src/users/route_user/gin_user"
)

func Routes(r *gin.RouterGroup, appCtx app_context.Appcontext) {
	r.POST("/login", gin_user.Login(appCtx))
	r.POST("/register", gin_user.Register(appCtx))
	authRoute := r.Group("/Private", middleware.RequireAuth(appCtx))
	{
		authRoute.GET("/infor", gin_user.GetProfile(appCtx))
		authRoute.GET("/address", gin_user.GetAddress(appCtx))
		authRoute.PUT("/update", gin_user.UpdateUser(appCtx))
		authRoute.POST("/adminUpdate", gin_user.AddUpdateAddmin(appCtx))
		// authRoute.POST("/register", userHandler.Register)
		// authRoute.POST("/login", userHandler.Login)
		// authRoute.POST("/refresh", refreshAuthMiddleware, userHandler.RefreshToken)
		// authRoute.GET("/me", authMiddleware, userHandler.GetMe)
		// authRoute.PUT("/change-password", authMiddleware, userHandler.ChangePassword)
	}
}
