package repository_orders

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	entities_orders "github.com/supersida159/e-commerce/src/order/entities_order"
)

func (s *sqlStore) GetOrder(ctx context.Context, id int) (*entities_orders.Order, error) {

	var data entities_orders.Order
	db := s.db
	if err := db.Table(data.TableName()).Where("id = ?", id).First(&data).Error; err != nil {
		return nil, common.ErrDB(err)
	}
	return &data, nil
}
