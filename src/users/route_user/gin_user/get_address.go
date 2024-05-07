package gin_user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/pkg/app_context"
	"github.com/supersida159/e-commerce/src/users/entities_user"
)

func GetAddress(appCtx app_context.Appcontext) gin.HandlerFunc {
	return func(c *gin.Context) {

		Data := c.MustGet(common.CurrentUser).(*entities_user.User)
		fmt.Println("Data ID", Data)

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(Data.Address))

	}
}
