package gin_order

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/common"
	response "github.com/supersida159/e-commerce/api-services/common/responese"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/pkg/kafka/saga"
)

func GetOrderStatusHandler(appCtx app_context.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("orderId")
		orchestrator := saga.NewOrchestrator(appCtx)

		orderEvent, err := orchestrator.GetOrderStatus(orderID)
		if err != nil {
			response.BuildErrorGinResponse(c, common.ErrOrderNotFound(err))
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(orderEvent))
	}
}
