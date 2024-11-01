package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
)

func Recover(ac app_context.Appcontext) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Header("Content-Type", "application/json")

				if appErr, ok := err.(*common.AppError); ok {
					c.AbortWithStatusJSON(appErr.StatusCode, appErr)
					return
				}
				appErr := common.ErrInternalServerError(err.(error))
				c.AbortWithStatusJSON(appErr.StatusCode, appErr)
				return
			}
		}()
		c.Next()
	}
}
