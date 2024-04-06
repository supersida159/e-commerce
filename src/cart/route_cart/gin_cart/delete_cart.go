package gin_carts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/pkg/app_context"
	repository_carts "github.com/supersida159/e-commerce/src/cart/repository_cart"
	usecase_carts "github.com/supersida159/e-commerce/src/cart/usecase_cart"
)

func DeleteCart(appCtx app_context.Appcontext) gin.HandlerFunc {
	return func(c *gin.Context) {
		userid := c.MustGet(common.CurrentUser).(common.Requester)
		store := repository_carts.NewSQLStore(appCtx.GetMainDBConnection())
		biz := usecase_carts.NewDeleteCartBiz(store)
		if err := biz.DeleteCartBiz(c.Request.Context(), userid.GetUserID()); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
