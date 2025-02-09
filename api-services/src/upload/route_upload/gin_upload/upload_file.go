package gin_upload

import (
	"errors"
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
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequestParameter(err))
			return
		}

		// Validate file size (e.g., limit to 10MB)
		const maxFileSize = 10 << 20 // 10MB
		if fileHeader.Size > maxFileSize {
			c.JSON(http.StatusRequestEntityTooLarge, common.ErrFailedToSavePhoto(errors.New("file size too large")))
			return
		}

		// Get folder name from form data (default is "img")
		folder := c.DefaultPostForm("folder", "img")

		// Open the uploaded file
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.ErrFailedToSavePhoto(err))
			return
		}
		defer file.Close()

		// Read file content into a byte slice
		dataBytes := make([]byte, fileHeader.Size)
		if _, err := file.Read(dataBytes); err != nil {
			c.JSON(http.StatusInternalServerError, common.ErrFailedToSavePhoto(err))
			return
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
				c.JSON(e.StatusCode, e)
			default:
				c.JSON(http.StatusInternalServerError, common.ErrInternalServerError(err))
			}
			return
		}

		// Return success response
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(img))
	}
}
