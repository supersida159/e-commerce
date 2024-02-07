package repositoy_upload

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	entities_Upload "github.com/supersida159/e-commerce/src/upload/entities_upload"
)

func (s *sqlStore) Create(ctx context.Context, data *entities_Upload.CreateUpload) error {
	db := s.db

	if err := db.Create(data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
