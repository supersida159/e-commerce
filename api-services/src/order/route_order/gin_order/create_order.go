package gin_order

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/common"
	response "github.com/supersida159/e-commerce/api-services/common/responese"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	repository_carts "github.com/supersida159/e-commerce/api-services/src/cart/repository_cart"
	"github.com/supersida159/e-commerce/api-services/src/order/DTO/order_request"
	"github.com/supersida159/e-commerce/api-services/src/order/DTO/order_response"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
	"github.com/supersida159/e-commerce/api-services/src/order/repository_orders"
	usecase_orders "github.com/supersida159/e-commerce/api-services/src/order/usecase_order"
)

func CreateOrderHandler(appCtx app_context.Appcontext) func(c *gin.Context) {

	return func(c *gin.Context) {
		var data order_response.OrderResponse
		var order entities_orders.Order
		var reqData order_request.CreateOrderRequest

		userContext := c.MustGet(common.CurrentUser).(common.Requester)

		storeOrder := repository_orders.NewSQLStore(appCtx.GetMainDBConnection())
		storeCart := repository_carts.NewSQLStore(appCtx.GetMainDBConnection())
		biz := usecase_orders.NewCreateOrderBiz(storeOrder,storeCart, appCtx.GetPubSub(), *appCtx.GetCache())

		if err := c.ShouldBindJSON(&reqData); err != nil {
			response.BuildErrorGinResponse(c, common.ErrJSONBlindding(err))
			return
		}
		validator := common.NewValidator()
		if appError := validator.ValidateStruct(reqData); appError != nil {
			response.BuildErrorGinResponse(c, appError)
			return
		}

		order = ConvertPlaceOrderReqToOrder(reqData)
		order.UserOrderID = userContext.GetUserID()

		if err := biz.CreateOrderBiz(c.Request.Context(), &order); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		order.Mask(true)

		//add convert to response

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.FakeId))
		response.BuildSuccessGinResponse(c, data.FakeId)
	}

}

func ConvertPlaceOrderReqToOrder(placeOrderReq order_request.CreateOrderRequest) entities_orders.Order {
	// Add appropriate initialization or mapping for UserOrderID, Shipping, OrderTotal, Notes, and OrderCancelled
	// Add any additional mapping or initialization logic for new fields
	return entities_orders.Order{
		CustomerName:   placeOrderReq.CustomerName,
		CustomerPhone:  placeOrderReq.CustomerPhone,
		CartID:         placeOrderReq.CartID,
		Shipping:       placeOrderReq.Shipping,
		Notes:          placeOrderReq.Notes,
		AddressID:      placeOrderReq.AddressID,
		OrderCancelled: false,
	}
}
