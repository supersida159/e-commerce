package generic_business

import (
	"context"

	"github.com/supersida159/e-commerce/api-services/common"
)

type Storage[T any] interface {
	FindById(ctx context.Context, Id int, moreInfo ...string) (*T, *common.AppError)
	Save(ctx context.Context, entity *T) (*T, *common.AppError)
	Delete(ctx context.Context, conditions map[string]interface{}) *common.AppError
	Update(ctx context.Context, updateData *T) (*T, *common.AppError)
	FindWithConditions(ctx context.Context, conditions map[string]interface{}, paging *common.Paging, orderClauses []string, moreInfo ...string) ([]T, *common.AppError)
}

type GenericsService[T any] struct {
	storage Storage[T]
}

func NewGenericsService[T any](storage Storage[T]) *GenericsService[T] {
	return &GenericsService[T]{storage: storage}
}

func (s *GenericsService[T]) FindById(ctx context.Context, Id int, moreInfo ...string) (*T, *common.AppError) {
	entity, err := s.storage.FindById(ctx, Id, moreInfo...)

	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *GenericsService[T]) Create(ctx context.Context, entity *T) (*T, *common.AppError) {
	// Save the entity using the storage layer
	newEntity, err := s.storage.Save(ctx, entity)
	if err != nil {
		return nil, err
	}
	return newEntity, nil
}

func (s *GenericsService[T]) Delete(ctx context.Context, conditions map[string]interface{}) *common.AppError {
	err := s.storage.Delete(ctx, conditions)
	if err != nil {
		return err
	}
	return nil
}

func (s *GenericsService[T]) Update(ctx context.Context, updateData *T) (*T, *common.AppError) {
	updatedEntity, err := s.storage.Update(ctx, updateData)
	if err != nil {
		return nil, err
	}
	return updatedEntity, nil
}

func (s *GenericsService[T]) FindList(ctx context.Context,
	conditions map[string]interface{},
	paging *common.Paging,
	orderClauses []string,
	moreInfo ...string) ([]T, *common.AppError) {

	entities, err := s.storage.FindWithConditions(ctx, conditions, paging, orderClauses, moreInfo...)
	if err != nil {
		return nil, err
	}
	return entities, nil
}
