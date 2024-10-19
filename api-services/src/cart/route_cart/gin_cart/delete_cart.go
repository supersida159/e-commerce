package gin_carts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/common"
	generic_business "github.com/supersida159/e-commerce/api-services/common/generics/business"
	generics_repository "github.com/supersida159/e-commerce/api-services/common/generics/repository"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	entities_carts "github.com/supersida159/e-commerce/api-services/src/cart/entities_cart"
)

func DeleteCart(appCtx app_context.Appcontext) gin.HandlerFunc {
	return func(c *gin.Context) {
		userid := c.MustGet(common.CurrentUser).(common.Requester)
		store := generics_repository.NewGenericStore[entities_carts.Cart](appCtx.GetMainDBConnection())
		biz := generic_business.NewGenericsService[entities_carts.Cart](store)
		conditions := map[string]interface{}{
			"UserID": userid.GetUserID(),
		}
		if err := biz.Delete(c.Request.Context(),conditions); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		} else {

		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
