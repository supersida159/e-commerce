package gin_order

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/src/order/DTO/order_request"
	"github.com/supersida159/e-commerce/api-services/src/order/DTO/order_response"
	"github.com/supersida159/e-commerce/api-services/src/order/repository_orders"
	usecase_orders "github.com/supersida159/e-commerce/api-services/src/order/usecase_order"
)

func CreateOrderHandler(appCtx app_context.Appcontext) func(c *gin.Context) {

	return func(c *gin.Context) {
		var data order_response.OrderResponse
		var reqData order_request.CreateOrderRequest

		userContext := c.MustGet(common.CurrentUser).(common.Requester)

		store := repository_orders.NewSQLStore(appCtx.GetMainDBConnection())
		biz := usecase_orders.NewCreateOrderBiz(store, appCtx.GetPubSub(), *appCtx.GetCache())

		if err := c.ShouldBind(&reqData); err != nil {
			if err == io.EOF {

			} else {
				c.JSON(http.StatusBadRequest, err)
				return
			}
		}

		data = ConvertPlaceOrderReqToOrder(reqData)
		data.UserOrderID = userContext.GetUserID()

		if err := biz.CreateOrderBiz(c.Request.Context(), &data); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		data.Mask(true)

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.FakeId))
	}

}

func ConvertPlaceOrderReqToOrder(placeOrderReq entities_orders.PlaceOrderReq) entities_orders.Order {
	order := entities_orders.Order{
		UserOrderID:    0,
		Shipping:       placeOrderReq.Shipping, // Add appropriate initialization or mapping
		OrderTotal:     0.0,                    // Add appropriate initialization or mapping
		Notes:          placeOrderReq.Notes,
		OrderCancelled: false, // Add appropriate initialization or mapping
	}

	// Add any additional mapping or initialization logic for new fields

	return order
}
