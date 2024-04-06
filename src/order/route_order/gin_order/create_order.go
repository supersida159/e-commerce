package gin_order

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/pkg/app_context"
	entities_orders "github.com/supersida159/e-commerce/src/order/entities_order"
	"github.com/supersida159/e-commerce/src/order/repository_orders"
	usecase_orders "github.com/supersida159/e-commerce/src/order/usecase_order"
)

func CreateOrderHandler(appCtx app_context.Appcontext) func(c *gin.Context) {

	return func(c *gin.Context) {
		var data entities_orders.Order
		var reqData entities_orders.PlaceOrderReq

		userContext := c.MustGet(common.CurrentUser).(common.Requester)

		store := repository_orders.NewSQLStore(appCtx.GetMainDBConnection())
		biz := usecase_orders.NewCreateOrderBiz(store, appCtx.GetPubSub(), *appCtx.GetCache())

		if err := c.ShouldBind(&reqData); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}
		if len(reqData.Products) == 0 {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(nil))
			return
		}
		data = ConvertPlaceOrderReqToOrder(reqData)
		data.UserOrderID = userContext.GetUserID()

		if err := biz.CreateOrderBiz(c.Request.Context(), &data); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		data.Mask(true)
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}

}

func ConvertPlaceOrderReqToOrder(placeOrderReq entities_orders.PlaceOrderReq) entities_orders.Order {
	order := entities_orders.Order{
		UserOrderID:    0,
		CustomerName:   placeOrderReq.CustomerName,
		CustomerPhone:  placeOrderReq.CustomerPhone,
		Shipping:       placeOrderReq.Shipping, // Add appropriate initialization or mapping
		OrderTotal:     0.0,                    // Add appropriate initialization or mapping
		Notes:          placeOrderReq.Notes,
		Address:        placeOrderReq.Address,
		OrderCancelled: false, // Add appropriate initialization or mapping
	}
	for _, p := range placeOrderReq.Products {

		order.Products = append(order.Products, &entities_orders.ProductQuantity{
			ProductID: p.ProductID,

			Quantity: p.Quantity,
		})
	}

	// Add any additional mapping or initialization logic for new fields

	return order
}
