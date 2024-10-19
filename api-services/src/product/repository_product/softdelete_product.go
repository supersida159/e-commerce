package repositoryproduct

import (
	"context"

	"github.com/supersida159/e-commerce/api-services/common"
)

func (s *sqlStore) SoftDeleteProduct(ctx context.Context, id int) error {
	db := s.db
	if err := db.Table("products").Where("id = ?", id).Updates(map[string]interface{}{"active": false}).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}
