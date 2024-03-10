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

func AddUpdateAddmin(appctx app_context.Appcontext) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := appctx.GetMainDBConnection()

		var data entities.UpdatePermission
		var updateUser entities.UserUpdate
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}
		admin := c.MustGet(common.CurrentUser).(common.Requester)

		if admin.GetRole() != "admin" {
			c.JSON(http.StatusUnauthorized, common.ErrInvalidRequest(nil))
			return
		}
		updateUser.Email = data.UpdateEmail
		updateUser.Role = data.RoleUpdate

		// missing check user submit the same data
		store := repository_user.NewSQLStore(db)
		md5 := hasher.NewMd5()
		repo := usecase_user.NewUpdatePermissionBusiness(appctx, store, md5)
		if err := repo.UpdateUserPermission(c.Request.Context(), &updateUser); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))

	}
}
