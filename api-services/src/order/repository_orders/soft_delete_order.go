package repository_orders

import (
	"context"

	"github.com/supersida159/e-commerce/api-services/common"
)

func (s *sqlStore) SoftDeleteOrder(ctx context.Context, id int) error {
	db := s.db
	if err := db.Table("orders").Where("id = ?", id).Updates(map[string]interface{}{"OrderCancelled": false, "Status": 0}).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}
