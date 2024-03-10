package ginproduct

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/pkg/app_context"
	"github.com/supersida159/e-commerce/src/product/entities_product"
	repositoryproduct "github.com/supersida159/e-commerce/src/product/repository_product"
	usecaseproduct "github.com/supersida159/e-commerce/src/product/usecase_product"
)

func CreateProductHandler(appCtx app_context.Appcontext) func(c *gin.Context) {

	return func(c *gin.Context) {
		var data entities_product.Product

		userPermission := c.MustGet(common.CurrentUser).(common.Requester)

		if userPermission.GetRole() != "admin" {
			c.JSON(http.StatusUnauthorized, common.ErrInvalidRequest(nil))
			return
		}
		store := repositoryproduct.NewSQLStore(appCtx.GetMainDBConnection())
		biz := usecaseproduct.NewCreateProductBiz(store)

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}
		if data.Category == "" {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(nil))
			return
		}
		data.Code = common.GenerateCode(data.Category)
		if err := biz.CreateProductBiz(c.Request.Context(), &data); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		data.Mask(true)
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}

}
