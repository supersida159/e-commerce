package gin_auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"golang.org/x/oauth2"
)

func OAuthHandler(appCtx app_context.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := appCtx.GetOAuth().AuthCodeURL("hello world", oauth2.AccessTypeOffline)
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}
