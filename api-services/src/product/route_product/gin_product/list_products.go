package ginproduct

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/src/product/entities_product"
	repositoryproduct "github.com/supersida159/e-commerce/api-services/src/product/repository_product"
	"github.com/supersida159/e-commerce/api-services/src/product/usecase_product"
)

func ListProducts(appCtx app_context.Appcontext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var reqData entities_product.ListProductReq
		var resData []entities_product.ListProductRes
		var paging common.Paging

		store := repositoryproduct.NewSQLStore(appCtx.GetMainDBConnection())
		biz := usecase_product.NewListProductsBiz(store)

		// if err := c.ShouldBind(&reqData); err != nil {
		// 	c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
		// 	return
		// }

		if err := c.ShouldBind(&reqData); err != nil {
			// Check if the error is due to an empty form
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
		category := c.Query("category")
		searchTerm := c.Query("searchTerm")
		if searchTerm != "" {
			reqData.SearchTerm = searchTerm
		}

		if category != "" {
			reqData.Category = category
		}
		paging.Fullfill()
		resData, err := biz.ListProductsBiz(c.Request.Context(), &reqData, &paging)
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
