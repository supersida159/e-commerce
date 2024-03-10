package usecase_orders

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	entities_orders "github.com/supersida159/e-commerce/src/order/entities_order"
)

type UpdateOrderStore interface {
	// FindOrder(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*entities_Order.Order, error)
	UpdateOrder(ctx context.Context, data *entities_orders.UpdateOrder) error
}

type updateOrderBiz struct {
	store UpdateOrderStore
}

func NewUpdateOrderBiz(store UpdateOrderStore) *updateOrderBiz {
	return &updateOrderBiz{
		store: store,
	}
}

func (biz *updateOrderBiz) UpdateOrderBiz(ctx context.Context, data *entities_orders.UpdateOrder) error {
	if err := biz.store.UpdateOrder(ctx, data); err != nil {
		return common.ErrCannotUpdateEntity(entities_orders.EntityName, err)
	}
	return nil
}
