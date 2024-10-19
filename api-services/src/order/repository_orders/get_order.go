package repository_orders

import (
	"context"

	"github.com/supersida159/e-commerce/api-services/common"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
)

func (s *sqlStore) GetOrder(ctx context.Context, id int) (*entities_orders.Order, error) {

	var data entities_orders.Order
	db := s.db
	if err := db.Table(data.TableName()).Where("id = ?", id).Preload("Cart.Items.Product").Preload("Address").First(&data).Error; err != nil {
		return nil, common.ErrDB(err)
	}
	return &data, nil
}
