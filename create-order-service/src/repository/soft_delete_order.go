package repository_orders

import (
	"context"

	"github.com/supersida159/e-commerce/create-order/common"
	entities "github.com/supersida159/e-commerce/create-order/src/model"
)

func (s *sqlStore) SoftDeleteOrder(ctx context.Context, id int) *common.AppError {
	db := s.db

	// Perform the soft delete using GORM's Delete method
	if err := db.Where("id = ?", id).Delete(&entities.Order{}).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
