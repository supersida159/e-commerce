package gin_order

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/pkg/app_context"
	"github.com/supersida159/e-commerce/src/product/entities_product"
	repositoryproduct "github.com/supersida159/e-commerce/src/product/repository_product"
	"github.com/supersida159/e-commerce/src/product/usecase_product"
)

func UpdateProductHandler(appCtx app_context.Appcontext) func(c *gin.Context) {
	return func(c *gin.Context) {

		userPermission := c.MustGet(common.CurrentUser).(common.Requester)
		if userPermission.GetRole() != "admin" {
			c.JSON(http.StatusUnauthorized, common.ErrInvalidRequest(nil))
			return
		}

		var data entities_product.Product
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}

		data.ID = int(data.FakeId.GetLocalID())
		store := repositoryproduct.NewSQLStore(appCtx.GetMainDBConnection())
		biz := usecase_product.NewUpdateProductBiz(store)
		if err := biz.UpdateProductBiz(c.Request.Context(), &data); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
