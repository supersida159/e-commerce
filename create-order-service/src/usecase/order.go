package usecase_orders

import (
	"context"

	"github.com/supersida159/e-commerce/create-order/common"
	entities "github.com/supersida159/e-commerce/create-order/src/model"
)

type OrderStore interface {
	CreateOrder(ctx context.Context, data *entities.Order) (*int, *common.AppError)
	SoftDeleteOrder(ctx context.Context, id int) *common.AppError
}

type OrderUsecase struct {
	store OrderStore
}

func NewOrderUsecase(store OrderStore) *OrderUsecase {
	return &OrderUsecase{store: store}
}

func (uc *OrderUsecase) CreateOrder(ctx context.Context, data *entities.Order) (*int, *common.AppError) {
	return uc.store.CreateOrder(ctx, data)
}

func (uc *OrderUsecase) SoftDeleteOrder(ctx context.Context, id int) *common.AppError {
	return uc.store.SoftDeleteOrder(ctx, id)
}
