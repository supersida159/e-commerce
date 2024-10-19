package gin_carts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	repository_carts "github.com/supersida159/e-commerce/api-services/src/cart/repository_cart"
	usecase_carts "github.com/supersida159/e-commerce/api-services/src/cart/usecase_cart"
)

func GetCart(appCtx app_context.Appcontext) gin.HandlerFunc {
	return func(c *gin.Context) {
		userid := c.MustGet(common.CurrentUser).(common.Requester)
		store := repository_carts.NewSQLStore(appCtx.GetMainDBConnection())
		biz := usecase_carts.NewGetCartBiz(store)
		cart, err := biz.GetCartBiz(c.Request.Context(), userid.GetUserID())
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		cart.Mask(true)
		for _, item := range cart.Items {
			item.Product.Mask(true)
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(cart))
	}
}
