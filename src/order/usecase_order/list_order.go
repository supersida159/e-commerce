package usecase_orders

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	entities_orders "github.com/supersida159/e-commerce/src/order/entities_order"
)

type ListOrdersStore interface {
	ListOrders(ctx context.Context,
		conditions map[string]interface{},
		filter *entities_orders.ListOrderReq,
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

func (biz *listOrdersBiz) ListOrdersBiz(ctx context.Context, fillter *entities_orders.ListOrderReq, paging *common.Paging, userid int) ([]entities_orders.Order, error) {
	var results []entities_orders.Order
	var err error
	if userid > 0 {
		results, err = biz.store.ListOrders(ctx, map[string]interface{}{"user_order_id": userid}, fillter, paging, "Cart", "Items", "Product")

	} else {
		results, err = biz.store.ListOrders(ctx, nil, fillter, paging, "Cart.Items.Product")

	}
	if err != nil {
		return nil, common.ErrCannotListEntity(entities_orders.EntityName, err)
	}
	return results, nil
}
