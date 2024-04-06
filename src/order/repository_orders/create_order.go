package repository_orders

import (
	"context"
	"fmt"

	"github.com/supersida159/e-commerce/common"
	entities_orders "github.com/supersida159/e-commerce/src/order/entities_order"
	"github.com/supersida159/e-commerce/src/product/entities_product"
)

func (s *sqlStore) CreateOrder(ctx context.Context, data *entities_orders.Order) error {
	db := s.db.Begin()
	for index, product := range data.Products {
		uid, _ := common.FromBase58(product.ProductID)
		productID := uid.GetLocalID()
		err := db.Table("products").Where("id = ?", productID).First(&data.Products[index].Product).Error
		if err != nil {
			return common.ErrDB(err)
		}
		if data.Products[index].Product.Quantity > data.Products[index].Quantity {
			if err := db.Table("products").Where("id = ?", productID).Update("quantity", data.Products[index].Product.Quantity-data.Products[index].Quantity).Error; err != nil {
				return common.ErrDB(err)
			}
		}
	}
	//check and update quantity product
	// for _, itemindb := range data.Products {
	// 	var item entities_product.Product
	// 	if err := db.Table("products").Where("name = ?", itemindb.Product.Name).First(&item).Error; err != nil {
	// 		return common.ErrDB(err)
	// 	}
	// 	if item.Quantity > itemindb.Quantity {
	// 		if err := db.Table("products").Where("name = ?", itemindb.Product.Name).Update("quantity", item.Quantity-itemindb.Quantity).Error; err != nil {
	// 			return common.ErrDB(err)
	// 		}
	// 	}
	// }
	//get total
	data.GetOrderTotal()
	//create an order

	if err := db.Table("orders").Create(&data).Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	return nil
}

func (s *sqlStore) OrderCancelled(ctx context.Context, data *entities_orders.Order) error {
	db := s.db.Begin()
	for _, itemindb := range data.Products {
		var item entities_product.Product
		if err := db.Table("products").Where("name = ?", itemindb.Product.Name).First(&item).Error; err != nil {
			return common.ErrDB(err)
		}
		if err := db.Table("products").Where("name = ?", itemindb.Product.Name).Update("quantity", item.Quantity+itemindb.Quantity).Error; err != nil {
			return common.ErrDB(err)
		}
	}
	if err := db.Table("orders").Where("id = ?", data.ID).Update("order_cancelled", true).Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}
	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)

	}
	fmt.Println("order_cancelled")
	return nil

}
