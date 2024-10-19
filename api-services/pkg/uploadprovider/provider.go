package uploadprovider

import (
	"context"

	"github.com/supersida159/e-commerce/api-services/common"
)

type UploadProvider interface {
	SaveFileUploaded(ctx context.Context, data []byte, dst string) (*common.Image, error)
}
