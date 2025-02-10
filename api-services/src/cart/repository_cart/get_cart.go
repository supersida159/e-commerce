package repository_carts

import (
	"context"
	"fmt"

	"github.com/supersida159/e-commerce/api-services/common"
	entities_carts "github.com/supersida159/e-commerce/api-services/src/cart/entities_cart"
)

func (s *sqlStore) GetCart(ctx context.Context, userID int, moreInfor ...string) (*entities_carts.Cart, *common.AppError) {
	db := s.db
	var cart entities_carts.Cart

	// Preload both CartItem and CartItem.Product
	if err := db.Preload("Items.Product").
		Where("UserID = ?", userID).
		Where("status = ?", 1).
		Where("deleted_at IS NULL").
		Order("id DESC").
		First(&cart).Error; err != nil {
		fmt.Println("Error fetching cart:", err)
		return nil, common.ErrDB(err)
	}

	return &cart, nil
}
