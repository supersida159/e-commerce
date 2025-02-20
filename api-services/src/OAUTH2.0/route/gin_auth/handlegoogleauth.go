package gin_auth

// import (
// 	"encoding/json"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/supersida159/e-commerce/api-services/common"
// 	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
// 	"github.com/supersida159/e-commerce/api-services/pkg/hasher"
// 	"github.com/supersida159/e-commerce/api-services/pkg/tokenprovider"
// 	"github.com/supersida159/e-commerce/api-services/pkg/tokenprovider/jwt"
// 	"github.com/supersida159/e-commerce/api-services/src/OAUTH2.0/entities_auth"
// 	usecase_oauth "github.com/supersida159/e-commerce/api-services/src/OAUTH2.0/usecase"
// 	"github.com/supersida159/e-commerce/api-services/src/users/entities_user"
// 	"github.com/supersida159/e-commerce/api-services/src/users/repository_user"
// 	"github.com/supersida159/e-commerce/api-services/src/users/usecase_user"
// 	"gorm.io/gorm"
// )

// func HandleGoogleOAuth(appCtx app_context.AppContext) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var req struct {
// 			AccessToken string `json:"access_token"`
// 		}

// 		// Parse JSON request body
// 		if err := c.ShouldBindJSON(&req); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
// 			return
// 		}

// 		// Verify the token with Google
// 		resp, err := http.Get("https://www.googleapis.com/oauth2/v3/userinfo?access_token=" + req.AccessToken)
// 		if err != nil || resp.StatusCode != http.StatusOK {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 			return
// 		}

// 		// Read and parse Google user info
// 		var googleUser entities_auth.GoogleUserInfo
// 		if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
// 			return
// 		}
// 		defer resp.Body.Close()

// 		// Check if user exists in your database
// 		store := repository_user.NewSQLStore(appCtx.GetMainDBConnection())
// 		tokenProvider := jwt.NewJwtProvider(appCtx.GetSecretKey())
// 		md5 := hasher.NewMd5()
// 		loginBussiness := usecase_user.NewLoginBusiness(appCtx, store, tokenProvider, md5)

// 		registerBussiness := usecase_user.NewRegisterBusiness(appCtx, store, md5)
// 		OAuthBusiness := usecase_oauth.NewOAuthBusiness(appCtx, loginBussiness, registerBussiness)
// 		existingUser, appErr := OAuthBusiness.StoreUser.FindUser(c.Request.Context(), map[string]interface{}{"email": googleUser.Email})
// 		if appErr != nil {
// 			if appErr.RootErr == gorm.ErrRecordNotFound {
// 				// Create new user
// 				newUser := entities_user.UserCreate{
// 					Email:     googleUser.Email,
// 					FirstName: googleUser.GivenName,
// 					LastName:  googleUser.FamilyName,
// 					Role:      "user",
// 					Avatar: &common.Image{
// 						Url: googleUser.Picture, // Assuming googleUser.Picture is the image URL from OAuth
// 					},
// 				}
// 				userDb, err := OAuthBusiness.StoreRegisterUser.RegisterGoogle(c.Request.Context(), &newUser)
// 				if err != nil {
// 					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
// 					return
// 				}
// 				userDb.Mask(false)
// 				payload := &tokenprovider.TokenPayload{
// 					UserId: userDb.ID,
// 					Role:   userDb.Role,
// 				}

// 				// Generate token for new user
// 				var account *entities_user.Account
// 				account, err = OAuthBusiness.StoreUser.GenerateToken(payload)
// 				if err != nil {
// 					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
// 					return
// 				}

// 				c.JSON(http.StatusOK, gin.H{
// 					"access_token": account,
// 				})
// 				return
// 			}

// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
// 			return
// 		} else {
// 			payload := &tokenprovider.TokenPayload{
// 				UserId: existingUser.ID,
// 				Role:   existingUser.Role,
// 			}
// 			// Generate token for new user
// 			var account *entities_user.Account
// 			account, appErr = OAuthBusiness.StoreUser.GenerateToken(payload)
// 			if appErr != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
// 				return
// 			}
// 			c.JSON(http.StatusOK, gin.H{
// 				"access_token": account.AccessToken,
// 			})
// 			return
// 		}
// 	}
// }
