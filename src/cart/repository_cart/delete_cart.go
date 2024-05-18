package repository_carts

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	entities_carts "github.com/supersida159/e-commerce/src/cart/entities_cart"
)

func (s *sqlStore) DeleteCart(ctx context.Context, userID int) error {
	db := s.db
	if err := db.Table(entities_carts.Cart{}.TableName()).Where("UserID = ?", userID).Where("status=?", 1).Update("status", 0).Error; err != nil {
		return common.ErrDB(err)
	}
	var newCart entities_carts.Cart
	newCart.UserID = userID
	if err := db.Table(entities_carts.Cart{}.TableName()).Create(&newCart).Error; err != nil {

	}
	return nil
}
