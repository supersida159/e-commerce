package usecase_carts

import (
	"context"

	"github.com/supersida159/e-commerce/api-services/common"
	generic_business "github.com/supersida159/e-commerce/api-services/common/generics/business"
	entities_carts "github.com/supersida159/e-commerce/api-services/src/cart/entities_cart"
	"github.com/supersida159/e-commerce/api-services/src/product/entities_product"
)

// ExtendedStorage extends the generic Storage interface with cart-specific methods
type ExtendedStorage interface {
	generic_business.Storage[entities_carts.Cart]
	UpdateCartItems(ctx context.Context, data *entities_product.CartItem, userID int) *common.AppError
}

// CartBiz embeds the generic business service
type CartBiz struct {
	*generic_business.GenericsService[entities_carts.Cart]
	store ExtendedStorage
}

// NewCartBiz creates a new CartBiz instance that extends GenericsService
func NewCartBiz(store ExtendedStorage) *CartBiz {
	return &CartBiz{
		GenericsService: generic_business.NewGenericsService[entities_carts.Cart](store),
		store:           store,
	}
}

func (s *CartBiz) UpdateCartItems(ctx context.Context, data *entities_product.CartItem, userID int) *common.AppError {
	if err := s.store.UpdateCartItems(ctx, data, userID); err != nil {
		return err
	}
	return nil
}

//	func (biz *CartBiz) CreateCartBiz(ctx context.Context, data *entities_carts.Cart) error {
//		if _, err := biz.storage.FindById(ctx, data.UserID); err == nil {
//			return common.ErrCannotCreateEntity(entities_carts.EntityName, err)
//		}
//		if err := biz.storage.Save(ctx, data); err != nil {
//			return common.ErrCannotCreateEntity(entities_carts.EntityName, err)
//		}
//		return nil
//	}
