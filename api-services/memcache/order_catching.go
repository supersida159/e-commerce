package memcache

import (
	"context"
	"fmt"

	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
)

func (c *memcaching) FindOrder(ctx context.Context, condition map[string]interface{}, moreInfores ...string) (*entities_orders.Order, error) {
	orderID := condition["id"].(int)

	orderInCache, err := c.store.Read(fmt.Sprintf("order-%d", orderID))
	if err != nil {
		return orderInCache.(*entities_orders.Order), nil
	}
	orderInRealStore, err := c.realStore.FinOrder(ctx, condition, moreInfores...)
	if err != nil {
		return nil, err
	}
	// go func() {
	// 	uc.store.Write(fmt.Sprintf("user-%d", userID), userInRealStore)
	// }()
	return orderInRealStore, nil

}
