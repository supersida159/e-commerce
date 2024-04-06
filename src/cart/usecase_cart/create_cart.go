package usecase_carts

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	entities_carts "github.com/supersida159/e-commerce/src/cart/entities_cart"
)

type CreateCartStore interface {
	CreateCart(ctx context.Context, data *entities_carts.Cart) error
	FindCart(ctx context.Context, userID int) (*entities_carts.Cart, error)
}

type createCartBiz struct {
	store CreateCartStore
}

func NewCreateCartBiz(store CreateCartStore) *createCartBiz {
	return &createCartBiz{
		store: store,
	}
}

func (biz *createCartBiz) CreateCartBiz(ctx context.Context, data *entities_carts.Cart) error {
	if _, err := biz.store.FindCart(ctx, data.UserID); err == nil {
		return common.ErrCannotCreateEntity(entities_carts.EntityName, err)
	}
	if err := biz.store.CreateCart(ctx, data); err != nil {
		return common.ErrCannotCreateEntity(entities_carts.EntityName, err)
	}
	return nil
}
