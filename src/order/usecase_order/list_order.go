package usecase_orders

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	entities_orders "github.com/supersida159/e-commerce/src/order/entities_order"
)

type ListOrdersStore interface {
	ListOrder(ctx context.Context,
		conditions map[string]interface{},
		paging *common.Paging,
		moreInfo ...string,
	) ([]entities_orders.Order, error)
}

type listOrdersBiz struct {
	store ListOrdersStore
}

func NewListOrdersBiz(store ListOrdersStore) *listOrdersBiz {
	return &listOrdersBiz{
		store: store,
	}
}

func (biz *listOrdersBiz) ListOrdersBiz(ctx context.Context, paging *common.Paging) ([]entities_orders.Order, error) {
	results, err := biz.store.ListOrder(ctx, nil, paging)
	if err != nil {
		return nil, common.ErrCannotListEntity(entities_orders.EntityName, err)
	}
	return results, nil
}
