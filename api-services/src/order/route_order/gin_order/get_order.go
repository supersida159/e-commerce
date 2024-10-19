package gin_order

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/src/order/repository_orders"
	usecase_orders "github.com/supersida159/e-commerce/api-services/src/order/usecase_order"
)

func Getorder(appCtx app_context.Appcontext) func(c *gin.Context) {
	return func(c *gin.Context) {
		store := repository_orders.NewSQLStore(appCtx.GetMainDBConnection())
		biz := usecase_orders.NewGetOrderBiz(store)
		req, _ := common.FromBase58(c.Param("id"))
		fmt.Println("Param:", c.Param("id"))
		fmt.Println("req:", req)
		fmt.Println("localID:", req.GetLocalID())

		data, err := biz.GetOrderBiz(c.Request.Context(), int(req.GetLocalID()))
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		data.Mask(true)
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}
}
