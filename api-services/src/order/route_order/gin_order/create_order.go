package gin_order

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/common"
	response "github.com/supersida159/e-commerce/api-services/common/responese"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/pkg/kafka/saga"
	"github.com/supersida159/e-commerce/api-services/src/order/DTO/order_request"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
	usecase_orders "github.com/supersida159/e-commerce/api-services/src/order/usecase_order"
)

func CreateOrderHandler(appCtx app_context.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			reqData order_request.CreateOrderRequest
			order   entities_orders.Order
		)

		userContext := c.MustGet(common.CurrentUser).(common.Requester)

		// Initialize the saga orchestrator
		producer := appCtx.GetProducer()
		orchestrator := saga.NewOrchestrator(producer, appCtx)

		// Initialize order business with proper configuration
		bizConfig := usecase_orders.OrderBusinessConfig{
			Producer:     producer,
			Orchestrator: orchestrator,
		}

		biz, err := usecase_orders.NewOrderBusiness(bizConfig)
		if err != nil {
			response.BuildErrorGinResponse(c, common.ErrInternalServerError(err))
			return
		}

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

		// Create order through the business layer
		if err := biz.CreateOrder(c.Request.Context(), &order); err != nil {
			response.BuildErrorGinResponse(c, err)
			return
		}

		// Mask sensitive data and prepare response
		order.Mask(true)

		// Return success response
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(&order))
	}
}

func ConvertPlaceOrderReqToOrder(placeOrderReq order_request.CreateOrderRequest) entities_orders.Order {
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
