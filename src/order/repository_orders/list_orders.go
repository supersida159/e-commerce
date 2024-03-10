package repository_orders

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	entities_orders "github.com/supersida159/e-commerce/src/order/entities_order"
)

func (s *sqlStore) ListOrders(ctx context.Context,
	conditions map[string]interface{},
	paging *common.Paging,
	moreInfo ...string,
) ([]entities_orders.Order, error) {
	var result []entities_orders.Order

	db := s.db.Table(entities_orders.Order{}.TableName()).Where(conditions)

	for _, info := range moreInfo {
		db = db.Preload(info)
	}

	// Find total count
	if err := db.Count(&paging.Total).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	// Apply pagination conditions

	if paging.FakeCusor != "" {
		uid, err := common.FromBase58(paging.FakeCusor)
		if err == nil {
			db = db.Where("id < ?", uid.GetLocalID())
		} else {
			db = db.Offset((paging.Page - 1) * paging.Limit)
		}
	}

	// Fetch the actual data
	if err := db.
		Limit(paging.Limit).
		Order("id desc").
		Find(&result).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	return result, nil
}
