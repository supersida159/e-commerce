package gin_upload

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/pkg/app_context"
	"github.com/supersida159/e-commerce/src/upload/repositoy_upload"
	"github.com/supersida159/e-commerce/src/upload/usecase_upload"
)

func UploadImg(appctx app_context.Appcontext) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := appctx.GetMainDBConnection()
		fileHeader, err := c.FormFile("file")
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}
		folder := c.DefaultPostForm("folder", "img")
		file, err := fileHeader.Open()

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		defer file.Close()

		dataBytes := make([]byte, fileHeader.Size)
		if _, err := file.Read(dataBytes); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		imgStore := repositoy_upload.NewSQLStore(db)
		biz := usecase_upload.NewUploadBiz(imgStore, appctx.UploadProvider(), appctx.GetPubSub())
		img, err := biz.Upload(c.Request.Context(), dataBytes, folder, fileHeader.Filename)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(img))

	}
}
