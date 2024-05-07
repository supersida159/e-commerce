package usecase_carts

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	entities_carts "github.com/supersida159/e-commerce/src/cart/entities_cart"
)

type GetCartStore interface {
	GetCart(ctx context.Context, userID int) (*entities_carts.Cart, error)
}

func NewGetCartBiz(store GetCartStore) *getCartBiz {
	return &getCartBiz{
		store: store,
	}
}

type getCartBiz struct {
	store GetCartStore
}

func (biz *getCartBiz) GetCartBiz(ctx context.Context, userID int) (*entities_carts.Cart, error) {

	cart, err := biz.store.GetCart(ctx, userID)
	if err != nil {
		return nil, common.ErrCannotGetEntity(entities_carts.EntityName, err)
	}
	return cart, nil
}
