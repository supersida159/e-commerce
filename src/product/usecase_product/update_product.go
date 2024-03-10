package usecase_product

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/src/product/entities_product"
)

type UpdateProductStore interface {
	// FindProduct(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*entities_product.Product, error)
	UpdateProduct(ctx context.Context, data *entities_product.Product) error
}

type updateProductBiz struct {
	store UpdateProductStore
}

func NewUpdateProductBiz(store UpdateProductStore) *updateProductBiz {
	return &updateProductBiz{
		store: store,
	}
}

func (biz *updateProductBiz) UpdateProductBiz(ctx context.Context, data *entities_product.Product) error {
	if err := biz.store.UpdateProduct(ctx, data); err != nil {
		return common.ErrCannotUpdateEntity(entities_product.EntityName, err)
	}
	return nil
}
