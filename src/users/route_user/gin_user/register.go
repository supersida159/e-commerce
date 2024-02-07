package gin_user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/pkg/app_context"
	"github.com/supersida159/e-commerce/pkg/hasher"
	"github.com/supersida159/e-commerce/src/users/entities"
	"github.com/supersida159/e-commerce/src/users/repository_user"
	"github.com/supersida159/e-commerce/src/users/usecase_user"
)

func Register(appctx app_context.Appcontext) func(c *gin.Context) {
	return func(c *gin.Context) {
		db := appctx.GetMainDBConnection()
		var data entities.UserCreate

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}
		store := repository_user.NewSQLStore(db)
		md5 := hasher.NewMd5()
		repo := usecase_user.NewRegisterBusiness(appctx, store, md5)

		if err := repo.Register(c.Request.Context(), &data); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		data.Mask(false)
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.FakeId.String()))
	}
}
