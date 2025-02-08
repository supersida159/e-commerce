package gin_upload

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/src/upload/repositoy_upload"
	"github.com/supersida159/e-commerce/api-services/src/upload/usecase_upload"
)

func UploadImg(appctx app_context.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := appctx.GetMainDBConnection()

		// Get file from form data
		fileHeader, err := c.FormFile("file")
		if err != nil {
			panic(common.ErrInvalidRequestParameter(err))
		}

		// Get folder name from form data (default is "img")
		folder := c.DefaultPostForm("folder", "img")

		// Open the uploaded file
		file, err := fileHeader.Open()
		if err != nil {
			panic(common.ErrFailedToSavePhoto(err))
		}
		defer file.Close()

		// Read file content into a byte slice
		dataBytes := make([]byte, fileHeader.Size)
		if _, err := file.Read(dataBytes); err != nil {
			panic(common.ErrFailedToSavePhoto(err))
		}

		// Initialize repository and business logic
		imgStore := repositoy_upload.NewSQLStore(db)
		biz := usecase_upload.NewUploadBiz(imgStore, appctx.UploadProvider(), appctx.GetPubSub())

		// Upload the image
		img, err := biz.Upload(c.Request.Context(), dataBytes, folder, fileHeader.Filename)
		if err != nil {
			// Handle specific errors from the business logic
			switch e := err.(type) {
			case *common.AppError:
				panic(e) // Re-throw AppError if it's already an AppError
			default:
				panic(common.ErrInternalServerError(err)) // Default to internal server error
			}
		}

		// Return success response
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(img))
	}
}
