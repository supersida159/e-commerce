package gin_carts

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/common"
	response "github.com/supersida159/e-commerce/api-services/common/responese"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	entities_carts "github.com/supersida159/e-commerce/api-services/src/cart/entities_cart"
	repository_carts "github.com/supersida159/e-commerce/api-services/src/cart/repository_cart"
	usecase_carts "github.com/supersida159/e-commerce/api-services/src/cart/usecase_cart"
	"github.com/supersida159/e-commerce/api-services/src/product/entities_product"
	"gorm.io/gorm"
)

type CartController struct {
	CartService *usecase_carts.CartBiz
	AppContext  app_context.Appcontext
	Validate    *common.Validator
}

func NewCartController(appContext app_context.Appcontext) (*CartController, *common.AppError) {
	dbs := appContext.GetMainDBConnection()

	store := repository_carts.NewCartStore(dbs)
	biz := usecase_carts.NewCartBiz(store)
	// service := generic_business.NewGenericsService[entities_carts.Cart](store)
	return &CartController{
		CartService: biz,
		AppContext:  appContext,
		Validate:    appContext.GetValidatetor(),
	}, nil
}

func (controller *CartController) CreateCart(c *gin.Context) {
	var data entities_carts.Cart
	userid := c.MustGet(common.CurrentUser).(common.Requester)

	// Check if a cart with the same UserID and status true already exists
	conditions := map[string]interface{}{
		"UserID": userid.GetUserID(),
		"Status": 1,
	}

	// check cart must be not exist (status = 0) before create newone
	_, err := controller.CartService.FindList(c.Request.Context(), conditions, nil, nil)
	if err != nil {
		if err.RootErr.Error() == gorm.ErrRecordNotFound.Error() {

		} else {
			response.BuildErrorGinResponseAndAbort(c, err)
			return
		}
	} else {
		response.BuildErrorGinResponseAndAbort(c, common.ErrDuplicateEntry(err))
		return
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		response.BuildErrorGinResponse(c, common.ErrJSONBlindding(err))
		return
	}
	data.UserID = userid.GetUserID()

	if responseData, err := controller.CartService.Create(c.Request.Context(), &data); err != nil {
		response.BuildErrorGinResponseAndAbort(c, err)
		return
	} else {
		response.BuildSuccessGinResponse(c, responseData)
	}
}

func (controller *CartController) DeleteCart(c *gin.Context) {
	userid := c.MustGet(common.CurrentUser).(common.Requester)

	var data entities_carts.Cart

	data.UserID = userid.GetUserID()

	data.Status = 0

	if _, err := controller.CartService.Update(c.Request.Context(), &data); err != nil {
		response.BuildErrorGinResponseAndAbort(c, err)
		return
	} else {
		_, err := controller.CartService.Create(c, &data)
		if err != nil {
			response.BuildErrorGinResponseAndAbort(c, err)
			return
		} else {
			response.BuildSuccessGinResponse(c, common.SimpleSuccessResponse(true))
		}
	}

}

func (controller *CartController) GetCart(c *gin.Context) {
	userid := c.MustGet(common.CurrentUser).(common.Requester)

	queryParams := c.Request.URL.Query()

	// Build conditions map
	// conditions := make(map[string]interface{})

	// for key, values := range queryParams {
	// 	if len(values) > 0 {
	// 		value := values[0]

	// 		// Skip pagination, sort, and preload parameters
	// 		if key == "page" || key == "limit" || key == "sort" || key == "preload" {
	// 			continue
	// 		}

	// 		// Check for comparison operators
	// 		switch {
	// 		case strings.HasPrefix(value, "<="):
	// 			conditions[key+" <="] = strings.TrimPrefix(value, "<=")
	// 		case strings.HasPrefix(value, ">="):
	// 			conditions[key+" >="] = strings.TrimPrefix(value, ">=")
	// 		case strings.HasPrefix(value, "<"):
	// 			conditions[key+" <"] = strings.TrimPrefix(value, "<")
	// 		case strings.HasPrefix(value, ">"):
	// 			conditions[key+" >"] = strings.TrimPrefix(value, ">")
	// 		case strings.Contains(value, ","):
	// 			// Handle IN clause
	// 			conditions[key] = strings.Split(value, ",")
	// 		default:
	// 			// Default to equality
	// 			conditions[key] = value
	// 		}
	// 	}
	// }

	// must add func for owner due to owener have many organization

	var paging common.Paging

	// must replace with filter later

	paging.Limit, _ = strconv.Atoi(c.Query("limit"))
	paging.Page, _ = strconv.Atoi(c.Query("page"))

	// Retrieve and parse sort parameters
	sortParams := queryParams["sort"]
	var orderClauses []string
	for _, sortParam := range sortParams {
		parts := strings.Split(sortParam, ",")
		if len(parts) == 2 {
			field, order := parts[0], strings.ToUpper(parts[1])
			if order != "ASC" && order != "DESC" {
				order = "ASC" // Default to ASC if the order is invalid
			}
			orderClauses = append(orderClauses, field+" "+order)
		}
	}

	// Retrieve and parse preload parameters
	preloadParams := queryParams["preload"]
	var preloadClauses []string
	for _, preloadParam := range preloadParams {
		preloadClauses = append(preloadClauses, strings.Split(preloadParam, ",")...)
	}

	//fill(set default) empty fielld
	paging.Fullfill()

	conditions := map[string]interface{}{
		"UserID": userid.GetUserID(),
		"Status": 1,
	}
	carts, err := controller.CartService.FindList(c.Request.Context(), conditions, &paging, orderClauses, preloadClauses...)
	if err != nil {
		response.BuildErrorGinResponseAndAbort(c, err)
		return
	}
	c.JSON(http.StatusOK, common.SimpleSuccessResponse(carts))
}

func (controller *CartController) UpdateCart(c *gin.Context) {
	var data entities_product.CartItem

	if err := c.ShouldBindJSON(&data); err != nil {
		response.BuildErrorGinResponse(c, common.ErrJSONBlindding(err))
		return
	}
	userid := c.MustGet(common.CurrentUser).(common.Requester)

	if err := controller.CartService.UpdateCartItemsBiz(c.Request.Context(), &data, userid.GetUserID()); err != nil {
		response.BuildErrorGinResponse(c, err)
		return
	}
	response.BuildSuccessGinResponse(c, data, "update Success")
}
