package usecase_carts

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	entities_carts "github.com/supersida159/e-commerce/src/cart/entities_cart"
)

type DeleteCartStore interface {
	DeleteCart(ctx context.Context, userID int) error
}

type deleteCartBiz struct {
	store DeleteCartStore
}

func NewDeleteCartBiz(store DeleteCartStore) *deleteCartBiz {
	return &deleteCartBiz{
		store: store,
	}
}
func (biz *deleteCartBiz) DeleteCartBiz(ctx context.Context, userID int) error {
	if err := biz.store.DeleteCart(ctx, userID); err != nil {
		return common.ErrCannotDeleteEntity(entities_carts.EntityName, err)
	}
	return nil
}
