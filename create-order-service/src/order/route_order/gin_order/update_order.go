package gin_order

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
	"github.com/supersida159/e-commerce/api-services/src/order/repository_orders"
	usecase_orders "github.com/supersida159/e-commerce/api-services/src/order/usecase_order"
)

func UpdateOrderHandler(appCtx app_context.Appcontext) func(c *gin.Context) {
	return func(c *gin.Context) {

		UserOrderID := c.MustGet(common.CurrentUser).(common.Requester)

		var data entities_orders.UpdateOrder
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}
		uidorderID, _ := common.FromBase58(c.Param("id"))
		orderID := int(uidorderID.GetLocalID())
		data.ID = orderID

		if data.Status != 1 {
			if UserOrderID.GetRole() != "admin" {
				c.JSON(http.StatusUnauthorized, common.ErrInvalidRequest(nil))
				return
			}
		}

		store := repository_orders.NewSQLStore(appCtx.GetMainDBConnection())
		biz := usecase_orders.NewUpdateOrderBiz(store)
		if err := biz.UpdateOrderBiz(c.Request.Context(), &data, UserOrderID.GetUserID()); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
