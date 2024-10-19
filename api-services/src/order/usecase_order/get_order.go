package usecase_orders

import (
	"context"

	"github.com/supersida159/e-commerce/api-services/common"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
)

type GetOrderStore interface {
	GetOrder(ctx context.Context, id int) (*entities_orders.Order, error)
}

type getOrderBiz struct {
	store GetOrderStore
}

func NewGetOrderBiz(store GetOrderStore) *getOrderBiz {
	return &getOrderBiz{store: store}
}
func (biz *getOrderBiz) GetOrderBiz(ctx context.Context, id int) (*entities_orders.Order, error) {
	data, err := biz.store.GetOrder(ctx, id)
	if err != nil {
		return nil, common.ErrCannotGetEntity(entities_orders.EntityName, err)
	}
	return data, nil
}
