package usecase_product

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/src/product/entities_product"
)

type CreateProductStore interface {
	CreateProduct(ctx context.Context, data *entities_product.Product) error
}

type createProductBiz struct {
	store CreateProductStore
}

func NewCreateProductBiz(store CreateProductStore) *createProductBiz {
	return &createProductBiz{
		store: store,
	}
}

func (biz *createProductBiz) CreateProductBiz(ctx context.Context, data *entities_product.Product) error {

	if err := biz.store.CreateProduct(ctx, data); err != nil {
		return common.ErrCannotCreateEntity(entities_product.EntityName, err)
	}
	return nil

}
