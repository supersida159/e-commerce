package repository_carts

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	entities_carts "github.com/supersida159/e-commerce/src/cart/entities_cart"
	"github.com/supersida159/e-commerce/src/product/entities_product"
	"gorm.io/gorm"
)

func (s *sqlStore) UpdateCart(ctx context.Context, data *entities_product.CartItem, userID int) error {
	var cart entities_carts.Cart
	err := s.db.Table(entities_carts.Cart{}.TableName()).Where("UserID = ?", userID).Preload("Items").First(&cart).Error
	if err == gorm.ErrRecordNotFound {
		//create cart
		err = s.CreateCart(ctx, &entities_carts.Cart{UserID: userID, Items: []*entities_product.CartItem{data}})
		if err != nil {
			return common.ErrDB(err)
		}
	} else if err != nil {
		// err := s.db.Table(entities_product.CartItem{}.TableName()).Where("UserID = ?", userID).First(&cart).Error

		return common.ErrDB(err)

	} else {
		//update cart

		for _, itemindb := range cart.Items {
			if itemindb.ProductID == data.ProductID {
				if data.Quantity == 0 {
					s.db.Table(entities_product.CartItem{}.TableName()).Where("id = ?", itemindb.ID).Delete(&itemindb)
					return nil
				} else {
					s.db.Table(entities_product.CartItem{}.TableName()).Where("id = ?", itemindb.ID).Update("quantity", data.Quantity)
					return nil
				}
			}

		}
		data.CartID = cart.ID
		s.db.Table(entities_product.CartItem{}.TableName()).Create(&data)

	}

	return nil
}
