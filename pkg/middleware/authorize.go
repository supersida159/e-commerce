package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/pkg/app_context"
	"github.com/supersida159/e-commerce/pkg/tokenprovider/jwt"
	"github.com/supersida159/e-commerce/src/users/entities"
)

type AuthenStore interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*entities.User, error)
}

func ErrInvalidToken() *common.AppError {
	return common.NewCustomError(
		errors.New("invalid token"),
		"invalid token",
		"ErrInvalidToken",
	)
}

func ExtractTokenFromHeaderString(c string) (string, error) {
	parts := strings.Split(c, " ")

	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", ErrInvalidToken()
	}
	return parts[1], nil
}

func RequireAuth(appCtx app_context.Appcontext) func(c *gin.Context) {
	tokenProvider := jwt.NewJwtProvider(appCtx.GetSecretKey())
	return func(c *gin.Context) {
		fmt.Println("Token:", c.GetHeader("Authorization"))
		token, err := ExtractTokenFromHeaderString(c.GetHeader("Authorization"))
		if err != nil {
			logrus.Fatal(err)
		}
		payload, err := tokenProvider.Validate(token)
		if err != nil {
			panic(err)
		}
		// db := appCtx.GetMainDBConnection()
		// store := usersstore.NewSQLStore(db)

		user, err := appCtx.GetCache().RealStore.FindUser(c.Request.Context(), map[string]interface{}{"id": payload.UserId})

		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
		}

		if user.Status == 0 {
			panic(common.ErrNoPermission(errors.New("user has been deleted of baned")))
		}
		user.Mask(false)
		c.Set(common.CurrentUser, user)
		c.Next()
	}

}