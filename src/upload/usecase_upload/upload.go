package usecase_upload

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/pkg/uploadprovider"
	entities_Upload "github.com/supersida159/e-commerce/src/upload/entities_upload"
)

type CreateImageStore interface {
	Create(ctx context.Context, data *entities_Upload.CreateUpload) error
}

type uploadBiz struct {
	provider uploadprovider.UploadProvider
	imgStore CreateImageStore
}

func NewUploadBiz(imgStore CreateImageStore, provider uploadprovider.UploadProvider) *uploadBiz {
	return &uploadBiz{imgStore: imgStore, provider: provider}
}

func (biz *uploadBiz) Upload(ctx context.Context, data []byte, folder, fileName string) (*common.Image, error) {
	fileBytes := bytes.NewReader(data)

	w, h, err := getImageDimension(fileBytes)
	if err != nil {
		return nil, entities_Upload.ErrCannotSaveFile(err)
	}
	if strings.TrimSpace(folder) == "" {
		folder = "img"
	}

	fileExt := filepath.Ext(fileName)                              // return .png, .jpg....
	fileName = fmt.Sprintf("%d%s", time.Now().UnixNano(), fileExt) //naming for file ex: 180283.png

	img, err := biz.provider.SaveFileUploaded(ctx, data, fmt.Sprintf("%s/%s", folder, fileName))
	if err != nil {
		return nil, entities_Upload.ErrCannotSaveFile(err)
	}

	img.Height = h
	img.Width = w
	img.CloudName = "s3"
	img.Extension = fileExt

	if err != nil {
		return nil, entities_Upload.ErrCannotSaveFile(err)
	}
	return img, nil

}

func getImageDimension(reader io.Reader) (int, int, error) {
	img, _, err := image.DecodeConfig(reader)
	if err != nil {
		return 0, 0, err
	}
	return img.Width, img.Height, nil

}
