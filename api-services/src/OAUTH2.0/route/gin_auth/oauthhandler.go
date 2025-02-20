// gin_auth/auth_handler.go
package gin_auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/pkg/hasher"
	"github.com/supersida159/e-commerce/api-services/pkg/tokenprovider"
	"github.com/supersida159/e-commerce/api-services/pkg/tokenprovider/jwt"
	dto "github.com/supersida159/e-commerce/api-services/src/OAUTH2.0/DTO"
	"github.com/supersida159/e-commerce/api-services/src/OAUTH2.0/services_auth"
	usecase_oauth "github.com/supersida159/e-commerce/api-services/src/OAUTH2.0/usecase"
	"github.com/supersida159/e-commerce/api-services/src/users/entities_user"
	"github.com/supersida159/e-commerce/api-services/src/users/repository_user"
	"github.com/supersida159/e-commerce/api-services/src/users/usecase_user"
	"gorm.io/gorm"
)

// HandleGoogleLogin handles the initial login request from Next.js
func HandleGoogleLogin(appCtx app_context.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		googleAuth, err := services_auth.NewGoogleAuthService()
		if err != nil {
			panic(err) // Handle this more gracefully in production
		}

		var userData dto.GoogleUserInfo
		if err := c.BindJSON(&userData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			return
		}

		// Verify the Google access token
		tokenInfo, err := googleAuth.VerifyGoogleToken(userData.AccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Google token"})
			return
		}
		// Verify that the email matches
		if tokenInfo.Email != userData.Email {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email mismatch"})
			return
		}
		// Initialize your stores and services
		store := repository_user.NewSQLStore(appCtx.GetMainDBConnection())
		tokenProvider := jwt.NewJwtProvider(appCtx.GetSecretKey())
		md5 := hasher.NewMd5()
		loginBusiness := usecase_user.NewLoginBusiness(appCtx, store, tokenProvider, md5)
		registerBusiness := usecase_user.NewRegisterBusiness(appCtx, store, md5)
		oAuthBusiness := usecase_oauth.NewOAuthBusiness(appCtx, loginBusiness, registerBusiness)

		// First, try to find the user
		user, appErr := loginBusiness.FindUser(c.Request.Context(), map[string]interface{}{
			"email": userData.Email,
		})

		var account *entities_user.Account

		if appErr != nil {
			// User doesn't exist, create new user
			// Create Image struct from Google profile picture
			// If it's not a "record not found" error, return the actual error
			if appErr.RootErr != gorm.ErrRecordNotFound {
				c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.RootErr.Error()})
				return
			}

			avatar := &common.Image{
				Url:       userData.Picture,
				Width:     500,      // Default width, adjust as needed
				Height:    500,      // Default height, adjust as needed
				CloudName: "google", // Indicate this is from Google
				Extension: "jpg",    // Default extension, adjust based on actual image
			}

			userCreate := &entities_user.UserCreate{
				Email:     userData.Email,
				FirstName: userData.Name,
				LastName:  userData.FamilyName,
				Password:  common.GenSalt(32), // Generate a random password for OAuth users
				Avatar:    avatar,
				Role:      "user",
				Salt:      common.GenSalt(50),
			}

			newUser, createErr := oAuthBusiness.StoreRegisterUser.RegisterGoogle(c.Request.Context(), userCreate)
			if createErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": createErr.Error()})
				return
			}

			// Generate token for new user
			payload := &tokenprovider.TokenPayload{
				UserId: newUser.ID,
				Role:   newUser.Role,
			}

			account, appErr = oAuthBusiness.StoreUser.GenerateToken(payload)
		} else {
			// User exists, update avatar if empty
			if user.Avatar == nil || user.Avatar.Url == "" {
				avatar := &common.Image{
					Url:       userData.Picture,
					Width:     500,      // Default width, adjust as needed
					Height:    500,      // Default height, adjust as needed
					CloudName: "google", // Indicate this is from Google
					Extension: "jpg",    // Default extension, adjust based on actual image
				}

				updateData := map[string]interface{}{
					"avatar": avatar,
				}
				if updateErr := store.UpdateUserAvatar(c.Request.Context(), user.ID, updateData); updateErr != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update avatar"})
					return
				}
			}

			// Generate token for existing user
			payload := &tokenprovider.TokenPayload{
				UserId: user.ID,
				Role:   user.Role,
			}

			account, appErr = oAuthBusiness.StoreUser.GenerateToken(payload)
		}

		if appErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, account)
	}
}

// AuthMiddleware protects routes requiring authentication
func AuthMiddleware(appCtx app_context.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("auth-token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No authentication token"})
			return
		}

		// Verify token using your token provider
		tokenProvider := jwt.NewJwtProvider(appCtx.GetSecretKey())
		claims, err := tokenProvider.Validate(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Set claims in context
		c.Set("user", claims)
		c.Next()
	}
}
