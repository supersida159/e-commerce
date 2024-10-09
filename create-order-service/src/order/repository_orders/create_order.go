package repository_orders

import (
	"context"
	"fmt"

	"github.com/supersida159/e-commerce/api-services/common"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
	"github.com/supersida159/e-commerce/api-services/src/product/entities_product"
)

func (s *sqlStore) CreateOrder(ctx context.Context, data *entities_orders.Order) error {
	db := s.db.Begin()
	err := db.Table("Cart").Preload("Items.Product").Where("UserID = ?", data.UserOrderID).Where("status = 1").Last(&data.Cart).Error
	if err != nil {
		return common.ErrCannotGetEntity(entities_orders.EntityName, err)
	}
	for index, product := range data.Cart.Items {
		err := db.Table("products").Where("id = ?", product.ProductID).First(&data.Cart.Items[index].Product).Error
		if err != nil {
			return common.ErrDB(err)
		}
		if data.Cart.Items[index].Product.Quantity > data.Cart.Items[index].Quantity {
			if err := db.Table("products").Where("id = ?", product.ProductID).Update("quantity", data.Cart.Items[index].Product.Quantity-data.Cart.Items[index].Quantity).Error; err != nil {
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
	data.Status = 2
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
	for _, itemindb := range data.Cart.Items {
		var item entities_product.Product
		if err := db.Table("products").Where("id = ?", itemindb.Product.ID).First(&item).Error; err != nil {
			return common.ErrDB(err)
		}
		if err := db.Table("products").Where("id = ?", itemindb.Product.ID).Update("quantity", item.Quantity+itemindb.Quantity).Error; err != nil {
			return common.ErrDB(err)
		}
	}
	if err := db.Table("orders").Where("id = ?", data.ID).Update("order_cancelled", true).Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}
	if err := db.Table("orders").Where("id = ?", data.ID).Update("status", 1).Error; err != nil {
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
