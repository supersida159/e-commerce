package usecase_orders

import (
	"context"
	"time"

	"github.com/supersida159/e-commerce/common"
	entities_orders "github.com/supersida159/e-commerce/src/order/entities_order"
)

type UpdateOrderStore interface {
	FindOrder(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*entities_orders.Order, error)
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

func (biz *updateOrderBiz) UpdateOrderBiz(ctx context.Context, data *entities_orders.UpdateOrder, userOderID int) error {
	previusOrder, err := biz.store.FindOrder(ctx, map[string]interface{}{"id": data.ID, "user_order_id": userOderID})

	if err != nil {
		return common.ErrCannotGetEntity("order", err)
	}
	if data.Address != nil {
		data.Address.UserID = userOderID
		if (previusOrder.Status == 2) && (previusOrder.Address == nil) {
			data.Shipping.EstimatedDelivery = time.Now().Add(3 * 24 * time.Hour)
			data.Shipping.Method = "COD"
			data.Status = 3
		}
	}

	if err := biz.store.UpdateOrder(ctx, data); err != nil {
		return common.ErrDB(err)
	}
	return nil
}
