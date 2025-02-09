package repository_orders

import (
	"context"

	"github.com/supersida159/e-commerce/create-order/common"
	entities "github.com/supersida159/e-commerce/create-order/src/model"
)

func (s *sqlStore) CreateOrder(ctx context.Context, data *entities.Order) (*int, *common.AppError) {

	if err := s.db.Table("orders").Create(&data).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	return &data.ID, nil
	// return common.ErrInternalServerError(errors.New("testing kafka Error"))
}

// func (s *sqlStore) OrderCancelled(ctx context.Context, data *entities.Order) error {
// 	db := s.db.Begin()
// 	for _, itemindb := range data.Cart.Items {
// 		var item entities_product.Product
// 		if err := db.Table("products").Where("id = ?", itemindb.Product.ID).First(&item).Error; err != nil {
// 			return common.ErrDB(err)
// 		}
// 		if err := db.Table("products").Where("id = ?", itemindb.Product.ID).Update("quantity", item.Quantity+itemindb.Quantity).Error; err != nil {
// 			return common.ErrDB(err)
// 		}
// 	}
// 	if err := db.Table("orders").Where("id = ?", data.ID).Update("order_cancelled", true).Error; err != nil {
// 		db.Rollback()
// 		return common.ErrDB(err)
// 	}
// 	if err := db.Table("orders").Where("id = ?", data.ID).Update("status", 1).Error; err != nil {
// 		db.Rollback()
// 		return common.ErrDB(err)
// 	}
// 	if err := db.Commit().Error; err != nil {
// 		db.Rollback()
// 		return common.ErrDB(err)

// 	}
// 	fmt.Println("order_cancelled")
// 	return nil

// }
