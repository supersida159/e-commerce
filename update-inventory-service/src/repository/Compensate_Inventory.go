package repository

import (
	"context"

	"github.com/supersida159/e-commerce/update-inventory-service/common"
	entities "github.com/supersida159/e-commerce/update-inventory-service/src/model"
	"gorm.io/gorm"
)

func (s *sqlStore) CompensateInventory(ctx context.Context, data entities.Cart) *common.AppError {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, item := range data.Items {
		var product entities.Product

		// Find the product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			return common.ErrEntityNotExist("product", err)
		}

		// Restore product quantity
		if err := tx.Model(&product).
			Update("quantity", gorm.Expr("quantity + ?", item.Quantity)).
			Error; err != nil {
			tx.Rollback()
			return common.ErrDB(err)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
