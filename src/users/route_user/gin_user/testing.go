package gin_user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/pkg/app_context"
)

func TestGetProfile(appctx app_context.Appcontext) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
