package gin_order

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/common"
	response "github.com/supersida159/e-commerce/api-services/common/responese"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	repository_carts "github.com/supersida159/e-commerce/api-services/src/cart/repository_cart"
	"github.com/supersida159/e-commerce/api-services/src/order/DTO/order_request"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
	"github.com/supersida159/e-commerce/api-services/src/order/repository_orders"
	usecase_orders "github.com/supersida159/e-commerce/api-services/src/order/usecase_order"
)

func CreateOrderHandler(appCtx app_context.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			reqData order_request.CreateOrderRequest
			order   entities_orders.Order
		)

		userContext := c.MustGet(common.CurrentUser).(common.Requester)

		// Bind and validate request data
		if err := c.ShouldBindJSON(&reqData); err != nil {
			response.BuildErrorGinResponse(c, common.ErrJSONBlindding(err))
			return
		}

		validator := common.NewValidator()
		if appError := validator.ValidateStruct(reqData); appError != nil {
			response.BuildErrorGinResponse(c, appError)
			return
		}

		// Convert request to order entity
		order = ConvertPlaceOrderReqToOrder(reqData)
		order.UserOrderID = userContext.GetUserID()

		orderStore := repository_orders.NewSQLStore(appCtx.GetMainDBConnection())
		cartStore := repository_carts.NewSQLStore(appCtx.GetMainDBConnection())
		biz := usecase_orders.NewOrderBusiness(orderStore, cartStore, appCtx.GetPubSub())

		// Create order through the business layer
		sagaID, err := biz.CreateOrder(c.Request.Context(), &order)
		if err != nil {
			response.BuildErrorGinResponse(c, err)
			return
		}

		// Return success response with polling and WebSocket URLs
		response := map[string]interface{}{
			"status":         "accepted",
			"message":        "Order processing started",
			"orderStatusUrl": fmt.Sprintf("/order-status/%s", *sagaID),
			"webSocketUrl":   fmt.Sprintf("/ws/order-status?orderId=%s", *sagaID),
		}
		c.JSON(http.StatusAccepted, common.SimpleSuccessResponse(response))
	}
}

func ConvertPlaceOrderReqToOrder(placeOrderReq order_request.CreateOrderRequest) entities_orders.Order {
	return entities_orders.Order{
		CustomerName:   placeOrderReq.CustomerName,
		CustomerPhone:  placeOrderReq.CustomerPhone,
		Shipping:       placeOrderReq.Shipping,
		Notes:          placeOrderReq.Notes,
		AddressID:      placeOrderReq.AddressID,
		OrderCancelled: false,
	}
}
