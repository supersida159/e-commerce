package repository_carts

import (
	"context"

	"github.com/supersida159/e-commerce/api-services/common"
	generics_repository "github.com/supersida159/e-commerce/api-services/common/generics/repository"
	entities_carts "github.com/supersida159/e-commerce/api-services/src/cart/entities_cart"
	"github.com/supersida159/e-commerce/api-services/src/product/entities_product"
	"gorm.io/gorm"
)

// func (s *sqlStore) UpdateCart(ctx context.Context, data *entities_product.CartItem, userID int) error {
// 	var cart entities_carts.Cart
// 	err := s.db.Table(entities_carts.Cart{}.TableName()).Where("UserID = ?", userID).Where("status=?", 1).Preload("Items").First(&cart).Error
// 	if err == gorm.ErrRecordNotFound {
// 		//create cart
// 		err = s.CreateCart(ctx, &entities_carts.Cart{UserID: userID, Items: []*entities_product.CartItem{data}})
// 		if err != nil {
// 			return common.ErrDB(err)
// 		}
// 	} else if err != nil {
// 		// err := s.db.Table(entities_product.CartItem{}.TableName()).Where("UserID = ?", userID).First(&cart).Error

// 		return common.ErrDB(err)

// 	} else {
// 		//update cart

// 		for _, itemindb := range cart.Items {
// 			if itemindb.ProductID == data.ProductID {
// 				if data.Quantity == 0 {
// 					s.db.Table(entities_product.CartItem{}.TableName()).Where("id = ?", itemindb.ID).Delete(&itemindb)
// 					return nil
// 				} else {
// 					s.db.Table(entities_product.CartItem{}.TableName()).Where("id = ?", itemindb.ID).Update("quantity", data.Quantity)
// 					return nil
// 				}
// 			}

// 		}
// 		data.CartID = cart.ID
// 		s.db.Table(entities_product.CartItem{}.TableName()).Create(&data)

// 	}

// 	return nil
// }

type CartStore struct {
	generics_repository.GenericStore[*entities_carts.Cart]
}

func NewCartStore(db *gorm.DB) *CartStore {
	return &CartStore{
		GenericStore: generics_repository.GenericStore[*entities_carts.Cart]{Db: db}, // Removed the '&'
	}
}

// UpdateCartItems is the additional function specific to CartStore
func (s *CartStore) UpdateCartItems(ctx context.Context, data *entities_product.CartItem, userID int) *common.AppError {
	var cart entities_carts.Cart
	err := s.Db.Table(entities_carts.Cart{}.TableName()).Where("UserID = ?", userID).Where("status=?", 1).Preload("Items").First(&cart).Error
	if err == gorm.ErrRecordNotFound {
		//create cart
		err = s.Db.Table(entities_carts.EntityName).Create(&entities_carts.Cart{UserID: userID, Items: []*entities_product.CartItem{data}}).Error
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
					s.Db.Table(entities_product.CartItem{}.TableName()).Where("id = ?", itemindb.ID).Delete(&itemindb)
					return nil
				} else {
					s.Db.Table(entities_product.CartItem{}.TableName()).Where("id = ?", itemindb.ID).Update("quantity", data.Quantity)
					return nil
				}
			}

		}
		data.CartID = cart.ID
		s.Db.Table(entities_product.CartItem{}.TableName()).Create(&data)

	}

	return nil
}
