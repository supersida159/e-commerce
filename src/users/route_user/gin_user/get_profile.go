package gin_user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/pkg/app_context"
)

func GetProfile(appCtx app_context.Appcontext) gin.HandlerFunc {
	return func(c *gin.Context) {

		Data := c.MustGet(common.CurrentUser).(common.Requester)

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(Data))

	}
}
