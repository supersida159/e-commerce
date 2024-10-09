package gin_order

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
	"github.com/supersida159/e-commerce/api-services/src/order/repository_orders"
	usecase_orders "github.com/supersida159/e-commerce/api-services/src/order/usecase_order"
)

func ListOrders(appCtx app_context.Appcontext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var reqData entities_orders.ListOrderReq
		var resData []entities_orders.Order
		var paging common.Paging
		var userID int

		store := repository_orders.NewSQLStore(appCtx.GetMainDBConnection())
		biz := usecase_orders.NewListOrdersBiz(store)

		CurrentUser := c.MustGet(common.CurrentUser).(common.Requester)

		if CurrentUser.GetRole() == "admin" {
			userID = 0
		} else if CurrentUser.GetRole() == "user" {
			userID = CurrentUser.GetUserID()
		} else {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(nil))
			return
		}

		if err := c.ShouldBind(&reqData); err != nil {
			if errors.Is(err, io.EOF) {
				// Form is empty
				// Handle the empty form case here
				fmt.Println("Form is empty")
			} else {
				c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
				return
			}
		}
		paging.Limit, _ = strconv.Atoi(c.Query("limit"))
		paging.Page, _ = strconv.Atoi(c.Query("page"))
		paging.FakeCusor = c.Query("cursor")
		paging.Fullfill()
		reqData.Mask(false)
		resData, err := biz.ListOrdersBiz(c.Request.Context(), &reqData, &paging, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		for index, _ := range resData {
			resData[index].Mask(true)
			if index == len(resData)-1 {
				paging.NextCursor = resData[index].FakeId.String()
			}
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(resData, &paging, &reqData))
	}
}
