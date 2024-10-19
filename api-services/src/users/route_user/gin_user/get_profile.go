package gin_user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
)

func GetProfile(appCtx app_context.Appcontext) gin.HandlerFunc {
	return func(c *gin.Context) {

		Data := c.MustGet(common.CurrentUser).(common.Requester)
		fmt.Println("Data ID", Data)

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(Data))

	}
}
