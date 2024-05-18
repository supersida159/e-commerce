package gin_order

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/pkg/app_context"
	entities_orders "github.com/supersida159/e-commerce/src/order/entities_order"
	"github.com/supersida159/e-commerce/src/order/repository_orders"
	usecase_orders "github.com/supersida159/e-commerce/src/order/usecase_order"
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
		data.Shipping.EstimatedDelivery = time.Now().Add(3 * 24 * time.Hour)
		data.Shipping.Method = "COD"
		data.Status = 2
		store := repository_orders.NewSQLStore(appCtx.GetMainDBConnection())
		biz := usecase_orders.NewUpdateOrderBiz(store)
		if err := biz.UpdateOrderBiz(c.Request.Context(), &data, UserOrderID.GetUserID()); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
