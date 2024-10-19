package usecase_product

import (
	"context"

	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/src/product/entities_product"
)

type GetProductStore interface {
	GetProduct(ctx context.Context, name string) (*[]entities_product.Product, error)
}

type getProductBiz struct {
	store GetProductStore
}

func NewGetProductBiz(store GetProductStore) *getProductBiz {
	return &getProductBiz{
		store: store,
	}
}

func (biz *getProductBiz) GetProductBiz(ctx context.Context, name string) (*[]entities_product.Product, error) {
	products, err := biz.store.GetProduct(ctx, name)
	if err != nil {
		return nil, common.ErrCannotGetEntity(entities_product.EntityName, err)
	}
	return products, nil
}
