package usecase

import (
	"context"

	"github.com/supersida159/e-commerce/update-cart-service/common"
)

type CartStore interface {
	SoftDeleteCart(ctx context.Context, id int) *common.AppError
	RecoveryCart(ctx context.Context, id int) *common.AppError
}

type CartUsecase struct {
	store CartStore
}

func NewCartUsecase(store CartStore) *CartUsecase {
	return &CartUsecase{store: store}
}

func (uc *CartUsecase) SoftDeleteCart(ctx context.Context, id int) *common.AppError {
	return uc.store.SoftDeleteCart(ctx, id)
}
func (uc *CartUsecase) RecoveryCart(ctx context.Context, id int) *common.AppError {
	return uc.store.RecoveryCart(ctx, id)
}
