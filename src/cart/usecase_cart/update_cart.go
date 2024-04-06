package usecase_carts

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	entities_carts "github.com/supersida159/e-commerce/src/cart/entities_cart"
)

type UpdateCartStore interface {
	UpdateCart(ctx context.Context, data *entities_carts.Cart) error
	FindCart(ctx context.Context, userID int) (*entities_carts.Cart, error)
}

type updateCartBiz struct {
	store UpdateCartStore
}

func NewUpdateCartBiz(store UpdateCartStore) *updateCartBiz {
	return &updateCartBiz{
		store: store,
	}
}

func (biz *updateCartBiz) UpdateCartBiz(ctx context.Context, data *entities_carts.Cart) error {
	_, err := biz.store.FindCart(ctx, data.UserID)
	if err != nil {
		return common.ErrCannotUpdateEntity(entities_carts.EntityName, err)
	}
	if err := biz.store.UpdateCart(ctx, data); err != nil {
		return common.ErrCannotUpdateEntity(entities_carts.EntityName, err)
	}
	return nil
}
