package usecase_orders

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	entities_orders "github.com/supersida159/e-commerce/src/order/entities_order"
)

type CreateOrderStore interface {
	CreateOrder(ctx context.Context, data *entities_orders.Order) error
}

type createOrderBiz struct {
	store CreateOrderStore
}

func NewCreateOrderBiz(store CreateOrderStore) *createOrderBiz {
	return &createOrderBiz{
		store: store,
	}
}

func (biz *createOrderBiz) CreateOrderBiz(ctx context.Context, data *entities_orders.Order) error {
	err := biz.store.CreateOrder(ctx, data)
	if err != nil {
		return common.ErrCannotCreateEntity(entities_orders.EntityName, err)
	}
	return nil
}
