package repository_carts

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	entities_carts "github.com/supersida159/e-commerce/src/cart/entities_cart"
)

func (s *sqlStore) UpdateCart(ctx context.Context, data *entities_carts.Cart) error {
	db := s.db
	if err := db.Table(entities_carts.Cart{}.TableName()).Where("UserID = ?", data.UserID).Updates(&data).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}
