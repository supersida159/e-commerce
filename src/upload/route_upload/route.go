package route_upload

import (
	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/pkg/app_context"
	"github.com/supersida159/e-commerce/pkg/middleware"
	"github.com/supersida159/e-commerce/src/upload/route_upload/gin_upload"
)

func Routes(r *gin.RouterGroup, appCtx app_context.Appcontext) {

	authRoute := r.Group("upload", middleware.RequireAuth(appCtx))
	{
		authRoute.POST("/addImage", gin_upload.UploadImg(appCtx))
		// authRoute.PUT("/update", gin_user.UpdateUser(appCtx))
		// authRoute.POST("/register", userHandler.Register)
		// authRoute.POST("/login", userHandler.Login)
		// authRoute.POST("/refresh", refreshAuthMiddleware, userHandler.RefreshToken)
		// authRoute.GET("/me", authMiddleware, userHandler.GetMe)
		// authRoute.PUT("/change-password", authMiddleware, userHandler.ChangePassword)
	}
}
