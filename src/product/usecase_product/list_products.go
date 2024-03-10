package usecase_product

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/src/product/entities_product"
)

type ListProductsStore interface {
	ListProduct(ctx context.Context,
		conditions map[string]interface{},
		filter *entities_product.ListProductReq,
		paging *common.Paging,
		moreInfo ...string,
	) ([]entities_product.ListProductRes, error)
}

type listProductsBiz struct {
	store ListProductsStore
}

func NewListProductsBiz(store ListProductsStore) *listProductsBiz {
	return &listProductsBiz{
		store: store,
	}
}

func (biz *listProductsBiz) ListProductsBiz(ctx context.Context, filter *entities_product.ListProductReq, paging *common.Paging) ([]entities_product.ListProductRes, error) {
	results, err := biz.store.ListProduct(ctx, nil, filter, paging)
	if err != nil {
		return nil, common.ErrCannotListEntity(entities_product.EntityName, err)
	}
	return results, nil
}
