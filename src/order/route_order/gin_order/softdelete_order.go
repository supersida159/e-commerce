package gin_order

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/pkg/app_context"
	"github.com/supersida159/e-commerce/src/order/repository_orders"
	usecase_orders "github.com/supersida159/e-commerce/src/order/usecase_order"
)

func SoftDeleteProductHandler(appCtx app_context.Appcontext) func(c *gin.Context) {

	return func(c *gin.Context) {
		getParam := c.Query("id")
		if getParam == "" {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(nil))
			return
		}
		FakeId, err := common.FromBase58(getParam)
		fmt.Println(FakeId.GetLocalID())
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}

		userPermission := c.MustGet(common.CurrentUser).(common.Requester)
		if userPermission.GetRole() != "admin" {
			c.JSON(http.StatusUnauthorized, common.ErrNoPermission(nil))
			return
		}
		store := repository_orders.NewSQLStore(appCtx.GetMainDBConnection())
		biz := usecase_orders.NewSoftDeleteOrderBiz(store)
		if err := biz.SoftDeleteOrderBiz(c.Request.Context(), int(FakeId.GetLocalID())); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}

}
