package repository_orders

import (
	"context"

	"github.com/supersida159/e-commerce/api-services/common"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
	"gorm.io/gorm"
)

func (s *sqlStore) ListOrders(ctx context.Context,
	conditions map[string]interface{},
	fillter *entities_orders.ListOrderReq,
	paging *common.Paging,
	moreInfo ...string,
) ([]entities_orders.Order, error) {
	var result []entities_orders.Order

	db := s.db.Table(entities_orders.Order{}.TableName()).Where(conditions)

	db = buildQuery(db, fillter)

	// Find total count
	if err := db.Model(entities_orders.Order{}).Count(&paging.Total).Error; err != nil {
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
	for _, info := range moreInfo {
		err := db.Preload(info).Error
		if err != nil {
			return nil, common.ErrDB(err)
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

func buildQuery(db *gorm.DB, filter *entities_orders.ListOrderReq) *gorm.DB {
	if filter.CustomerName != "" {
		db = db.Where("customer_name = ?", filter.CustomerName)
	}
	if filter.CustomerPhone != "" {
		db = db.Where("customer_phone = ?", filter.CustomerPhone)
	}
	if filter.Status != 0 {
		db = db.Where("status = ?", filter.Status)
	}
	if filter.OrderCancelled {
		db = db.Where("order_cancelled = ?", filter.OrderCancelled)
	}
	if filter.CreatedAt != nil {
		db = db.Where("created_at > ?", filter.CreatedAt)
	}
	if filter.UpdatedAt != nil {
		db = db.Where("updated_at > ?", filter.UpdatedAt)
	}
	if filter.ID != 0 {
		db = db.Where("id = ?", filter.ID)
	}
	return db
}
