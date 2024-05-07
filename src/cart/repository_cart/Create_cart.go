package repository_carts

import (
	"context"
	"errors"

	"github.com/supersida159/e-commerce/common"
	entities_carts "github.com/supersida159/e-commerce/src/cart/entities_cart"
	"github.com/supersida159/e-commerce/src/product/entities_product"
	"gorm.io/gorm"
)

func (s *sqlStore) CreateCart(ctx context.Context, data *entities_carts.Cart) error {
	db := s.db
	//check and update quantity product
	for index, itemindb := range data.Items {
		var item entities_product.Product
		if err := db.Table("products").Where("id = ?", itemindb.ProductID).First(&item).Error; err == nil {
			if item.Quantity > itemindb.Quantity { // check if quantity is enough
				// if err := db.Table("products").Where("name  = ?", itemindb.Product.Name).Update("quantity", item.Quantity-itemindb.Quantity).Error; err != nil {
				// 	return common.ErrDB(err)
				// }
				data.Items[index].Quantity = itemindb.Quantity
			}
		} else {
			if len(data.Items) > 1 {
				data.Items = append(data.Items[:index], data.Items[index+1:]...)
			} else {
				return common.NewCustomError(errors.New("product not found"), "product not found", "ErrProductNotFound")
			}
		}

	}
	//create an order
	if err := db.Table("Cart").Create(&data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}

func (s *sqlStore) FindCart(ctx context.Context, userID int) (*entities_carts.Cart, error) {
	db := s.db.Table(entities_carts.Cart{}.TableName())

	var cart entities_carts.Cart

	if err := db.Where("UserID = ?", userID).First(&cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}

	return &cart, nil
}
