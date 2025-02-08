package usecase

import (
	"context"

	entities "github.com/supersida159/e-commerce/update-inventory-service/src/model"

	"github.com/supersida159/e-commerce/update-inventory-service/common"
)

type InventoryStore interface {
	UpdateInventory(ctx context.Context, data entities.Cart) *common.AppError
	CompensateInventory(ctx context.Context, data entities.Cart) *common.AppError
}

type InventoryUsecase struct {
	store InventoryStore
}

func NewInventoryUsecase(store InventoryStore) *InventoryUsecase {
	return &InventoryUsecase{store: store}
}

func (uc *InventoryUsecase) UpdateInventory(ctx context.Context, data entities.Cart) *common.AppError {
	return uc.store.UpdateInventory(ctx, data)
}
func (uc *InventoryUsecase) CompensateInventory(ctx context.Context, data entities.Cart) *common.AppError {
	return uc.store.CompensateInventory(ctx, data)
}
