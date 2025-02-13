package gin_auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
)

func OAuthCallbackHandler(appCtx app_context.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extracting the "code" query parameter from the request.
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing code parameter"})
			return
		}

		// Exchanging the code for an access token.
		token, err := appCtx.GetOAuth().Exchange(context.Background(), code)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Creating an HTTP client to make authenticated requests using the access token.
		client := appCtx.GetOAuth().Client(context.Background(), token)

		// Fetching the user's public details from the Google API endpoint.
		resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		// Decoding the JSON response into a generic map.
		var userInfo map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&userInfo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Sending the user's public details as a JSON response.
		c.JSON(http.StatusOK, userInfo)
	}
}
