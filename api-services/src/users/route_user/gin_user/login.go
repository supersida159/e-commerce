package gin_user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/pkg/hasher"
	"github.com/supersida159/e-commerce/api-services/pkg/tokenprovider/jwt"
	"github.com/supersida159/e-commerce/api-services/src/users/entities_user"
	"github.com/supersida159/e-commerce/api-services/src/users/repository_user"
	"github.com/supersida159/e-commerce/api-services/src/users/usecase_user"
)

func Login(appCtx app_context.Appcontext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var LoginUserData entities_user.UserLogin
		if err := c.ShouldBind(&LoginUserData); err != nil {
			panic(common.ErrInvalidRequest(err))
		}
		fmt.Println("LoginUserData", LoginUserData)
		db := appCtx.GetMainDBConnection()
		tokenProvider := jwt.NewJwtProvider(appCtx.GetSecretKey())

		store := repository_user.NewSQLStore(db)
		md5 := hasher.NewMd5()

		business := usecase_user.NewLoginBusiness(appCtx, store, tokenProvider, md5, 60*60*24*30)
		acount, err := business.Login(c.Request.Context(), &LoginUserData)

		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(acount))
	}
}
