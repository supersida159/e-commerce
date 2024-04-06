package ginproduct

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/pkg/app_context"
	repositoryproduct "github.com/supersida159/e-commerce/src/product/repository_product"
	"github.com/supersida159/e-commerce/src/product/usecase_product"
)

func GetProductHandler(appCtx app_context.Appcontext) func(c *gin.Context) {
	return func(c *gin.Context) {
		store := repositoryproduct.NewSQLStore(appCtx.GetMainDBConnection())
		biz := usecase_product.NewGetProductBiz(store)
		name := c.Param("name")

		data, err := biz.GetProductBiz(c.Request.Context(), name)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		for i := range *data {
			(*data)[i].Mask(true)
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}
}
