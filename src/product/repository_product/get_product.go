package repositoryproduct

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/src/product/entities_product"
)

func (s *sqlStore) GetProduct(ctx context.Context, name string) (*[]entities_product.Product, error) {
	var result []entities_product.Product
	if err := s.db.Table(entities_product.Product{}.TableName()).Where("name = ?", name).Find(&result).Error; err != nil {
		return nil, common.ErrDB(err)
	}
	return &result, nil
}
