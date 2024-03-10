package repositoryproduct

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/src/product/entities_product"
)

func (s *sqlStore) UpdateProduct(ctx context.Context, data *entities_product.Product) error {
	db := s.db

	if err := db.Table(data.TableName()).Where("id = ?", data.ID).Updates(&data).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}

// func (s *sqlStore) FindProduct(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*entities_product.Product, error) {
// 	db := s.db
// 	var product entities_product.Product
// 	if err := db.Where(conditions).First(&product).Error; err != nil {
// 		return nil, common.ErrDB(err)
// 	}
// 	return &product, nil
// }
