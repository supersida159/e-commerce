package repository_orders

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	entities_orders "github.com/supersida159/e-commerce/src/order/entities_order"
)

func (s *sqlStore) UpdateOrder(ctx context.Context, data *entities_orders.UpdateOrder) error {
	db := s.db
	if data.AddressID != 0 && data.Address != nil {
		err := db.Table("addresses").Where("id = ?", data.AddressID).Updates(&data.Address).Error
		if err != nil {
			return common.ErrCannotUpdateEntity("address", err)
		}
	} else {
		if data.Address != nil {
			err := db.Table("addresses").Create(&data.Address).Error
			if err != nil {

				return common.ErrCannotCreateEntity("address", err)
			}
		}
	}

	if err := db.Table("orders").Where("id = ?", data.ID).Updates(data).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}

func (s *sqlStore) FindOrder(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*entities_orders.Order, error) {
	db := s.db
	var data entities_orders.Order
	if err := db.Table(data.TableName()).Where(conditions).First(&data).Error; err != nil {
		return nil, common.ErrDB(err)
	}
	return &data, nil
}
