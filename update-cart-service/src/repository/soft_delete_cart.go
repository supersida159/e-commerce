package repository

import (
	"context"

	"github.com/supersida159/e-commerce/update-cart-service/common"
	entities "github.com/supersida159/e-commerce/update-cart-service/src/model"
)

func (s *sqlStore) SoftDeleteCart(ctx context.Context, id int) *common.AppError {
	db := s.db

	// Perform the soft delete using GORM's Delete method
	if err := db.Where("id = ?", id).Delete(&entities.Cart{}).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
