package gin_carts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	repository_carts "github.com/supersida159/e-commerce/api-services/src/cart/repository_cart"
	usecase_carts "github.com/supersida159/e-commerce/api-services/src/cart/usecase_cart"
	"github.com/supersida159/e-commerce/api-services/src/product/entities_product"
)

func UpdateCart(appCtx app_context.Appcontext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data entities_product.CartItem
		userid := c.MustGet(common.CurrentUser).(common.Requester)
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}

		if data.Product.FakeId != nil {
			uid, err := common.FromBase58(data.Product.FakeId.String())
			if err != nil {
				c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
				return
			} else {
				data.Product.ID = int(uid.GetLocalID())
				data.ProductID = data.Product.ID
			}

		}
		store := repository_carts.NewSQLStore(appCtx.GetMainDBConnection())
		biz := usecase_carts.NewUpdateCartBiz(store)
		if err := biz.UpdateCartBiz(c.Request.Context(), &data, userid.GetUserID()); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
