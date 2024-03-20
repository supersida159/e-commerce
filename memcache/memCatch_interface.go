package memcache

import (
	"context"

	entities_orders "github.com/supersida159/e-commerce/src/order/entities_order"
	"github.com/supersida159/e-commerce/src/users/entities_user"
)

type memcaching struct {
	store     Caching
	realStore RealStore
}

func NewMemCaching(store Caching, realStore RealStore) *memcaching {
	return &memcaching{
		store:     store,
		realStore: realStore,
	}
}

type RealStore interface {
	FindUser(ctx context.Context, condition map[string]interface{}, moreInfores ...string) (*entities_user.User, error)
	FinOrder(ctx context.Context, condition map[string]interface{}, moreInfores ...string) (*entities_orders.Order, error)
}
