package repository_carts

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	entities_carts "github.com/supersida159/e-commerce/src/cart/entities_cart"
)

func (s *sqlStore) GetCart(ctx context.Context, userID int) (*entities_carts.Cart, error) {
	db := s.db
	var cart entities_carts.Cart
	if err := db.Table(entities_carts.Cart{}.TableName()).Where("UserID = ?", userID).Preload("Items.Product").First(&cart).Error; err != nil {
		return nil, common.ErrDB(err)
	}
	return &cart, nil
}
