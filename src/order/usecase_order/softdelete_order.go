package usecase_orders

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	entities_orders "github.com/supersida159/e-commerce/src/order/entities_order"
)

type SoftDeleteOrderStore interface {
	// FindOrder(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*entities_Order.Order, error)
	SoftDeleteOrder(ctx context.Context, id int) error
}

type softDeleteOrderBiz struct {
	store SoftDeleteOrderStore
}

func NewSoftDeleteOrderBiz(store SoftDeleteOrderStore) *softDeleteOrderBiz {
	return &softDeleteOrderBiz{
		store: store,
	}
}

func (biz *softDeleteOrderBiz) SoftDeleteOrderBiz(ctx context.Context, id int) error {

	if err := biz.store.SoftDeleteOrder(ctx, id); err != nil {
		return common.ErrCannotDeleteEntity(entities_orders.EntityName, err)
	}
	return nil
}
