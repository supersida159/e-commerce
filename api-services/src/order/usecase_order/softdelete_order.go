package usecase_orders

import (
	"context"

	"github.com/supersida159/e-commerce/api-services/common"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
)

type SoftDeleteOrderStore interface {
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
