package repository_orders

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	entities_orders "github.com/supersida159/e-commerce/src/order/entities_order"
)

func (s *sqlStore) UpdateOrder(ctx context.Context, data *entities_orders.Order) error {
	db := s.db

	if err := db.Table(data.TableName()).Where("id = ?", data.ID).Updates(&data).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}
