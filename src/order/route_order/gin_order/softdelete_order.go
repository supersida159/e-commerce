package gin_order

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/pkg/app_context"
	repositoryproduct "github.com/supersida159/e-commerce/src/product/repository_product"
	"github.com/supersida159/e-commerce/src/product/usecase_product"
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
		store := repositoryproduct.NewSQLStore(appCtx.GetMainDBConnection())
		biz := usecase_product.NewSoftDeleteProductBiz(store)
		if err := biz.SoftDeleteProductBiz(c.Request.Context(), int(FakeId.GetLocalID())); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}

}
