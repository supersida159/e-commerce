package usecase_product

import (
	"context"

	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/src/product/entities_product"
)

type SoftDeleteProductStore interface {
	// FindProduct(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*entities_product.Product, error)
	SoftDeleteProduct(ctx context.Context, id int) error
}

type softDeleteProductBiz struct {
	store SoftDeleteProductStore
}

func NewSoftDeleteProductBiz(store SoftDeleteProductStore) *softDeleteProductBiz {
	return &softDeleteProductBiz{
		store: store,
	}
}

func (biz *softDeleteProductBiz) SoftDeleteProductBiz(ctx context.Context, id int) error {

	if err := biz.store.SoftDeleteProduct(ctx, id); err != nil {
		return common.ErrCannotDeleteEntity(entities_product.EntityName, err)
	}
	return nil
}
