package usecase_carts

import (
	"context"

	"github.com/supersida159/e-commerce/api-services/common"
	entities_carts "github.com/supersida159/e-commerce/api-services/src/cart/entities_cart"
	repository_carts "github.com/supersida159/e-commerce/api-services/src/cart/repository_cart"
)

type Storage[T any] interface {
	FindById(ctx context.Context, Id int, moreInfo ...string) (*T, *common.AppError)
	Save(ctx context.Context, entity *T) (*T, *common.AppError)
	Delete(ctx context.Context, id int) *common.AppError
	Update(ctx context.Context, updateData *T) (*T, *common.AppError)
	FindWithConditions(ctx context.Context, conditions map[string]interface{}, paging *common.Paging, orderClauses []string, moreInfo ...string) ([]T, *common.AppError)
}

type CartBiz struct {
	storage repository_carts.CartStore
}

func NewCartBiz(storage repository_carts.CartStore) *CartBiz {
	return &CartBiz{storage: storage}
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
func (s *CartBiz) Create(ctx context.Context, entity *entities_carts.Cart) (*entities_carts.Cart, *common.AppError) {

	// Save the entity using the storage layer
	newEntity, err := s.storage.Save(ctx, entity)
	if err != nil {
		return nil, err
	}
	return newEntity, nil
}

func (s *CartBiz) Delete(ctx context.Context, id int) *common.AppError {
	err := s.storage.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *CartBiz) Update(ctx context.Context, updateData *entities_carts.Cart) (*entities_carts.Cart, *common.AppError) {
	updatedEntity, err := s.storage.Update(ctx, updateData)
	if err != nil {
		return nil, err
	}
	return updatedEntity, nil
}

func (s *CartBiz) FindList(ctx context.Context,
	conditions map[string]interface{},
	paging *common.Paging,
	orderClauses []string,
	moreInfo ...string) ([]entities_carts.Cart, *common.AppError) {

	entities, err := s.storage.FindWithConditions(ctx, conditions, paging, orderClauses, moreInfo...)
	if err != nil {
		return nil, err
	}
	return entities, nil
}
