package gin_carts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/pkg/app_context"
	entities_carts "github.com/supersida159/e-commerce/src/cart/entities_cart"
	repository_carts "github.com/supersida159/e-commerce/src/cart/repository_cart"
	usecase_carts "github.com/supersida159/e-commerce/src/cart/usecase_cart"
)

func CreateCart(appCtx app_context.Appcontext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data entities_carts.Cart
		userid := c.MustGet(common.CurrentUser).(common.Requester)
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}
		data.UserID = userid.GetUserID()
		store := repository_carts.NewSQLStore(appCtx.GetMainDBConnection())
		biz := usecase_carts.NewCreateCartBiz(store)
		if err := biz.CreateCartBiz(c.Request.Context(), &data); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}
}
