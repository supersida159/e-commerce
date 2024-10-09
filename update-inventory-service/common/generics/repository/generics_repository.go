package generics_repository

import (
	"context"
	"fmt"
	"reflect"

	"github.com/supersida159/e-commerce/api-services/common"
	"gorm.io/gorm"
)

// GenericStore is a generic store that can work with any Entity
type GenericStore[T any] struct {
	db *gorm.DB
}

// NewGenericStore creates a new GenericStore
func NewGenericStore[T any](db *gorm.DB) *GenericStore[T] {
	return &GenericStore[T]{db: db}
}

// FindById is a generic function to find an entity by ID
func (s *GenericStore[T]) FindById(ctx context.Context, Id string, moreInfo ...string) (*T, error) {
	db := s.db.Begin()
	var entity T

	var tableName string

	// Check if T has TableName method
	if tableNamer, ok := any(new(T)).(interface{ TableName() string }); ok {
		tableName = tableNamer.TableName()
	} else {
		// Fallback to default table name if TableName method is not implemented
		entityType := reflect.TypeOf(new(T)).Elem().Name()
		tableName = s.db.NamingStrategy.TableName(entityType)
	}

	db = db.Table(tableName)

	for _, info := range moreInfo {
		db = db.Preload(info)
	}

	if err := db.WithContext(ctx).Where("id=?", Id).First(&entity).Error; err != nil {
		db.Rollback()
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrEntityNotExist(tableName, err)
		}
		return nil, common.ErrDB(err)
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return nil, common.ErrDB(err)
	}

	return &entity, nil
}

func (s *GenericStore[T]) Save(ctx context.Context, entity *T) (*T, error) {
	var tableName string

	// Check if T has TableName method
	if tableNamer, ok := any(new(T)).(interface{ TableName() string }); ok {
		tableName = tableNamer.TableName()
	} else {
		// Fallback to default table name if TableName method is not implemented
		entityType := reflect.TypeOf(new(T)).Elem().Name()
		tableName = s.db.NamingStrategy.TableName(entityType)
	}

	// Start a transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, common.ErrDB(tx.Error)
	}

	// Create the entity
	if err := tx.Table(tableName).Create(entity).Error; err != nil {
		tx.Rollback()
		return nil, common.ErrDB(err)
	}

	// Retrieve the primary key field
	primaryKey := tx.Statement.Schema.PrioritizedPrimaryField.DBName

	// Fetch the newly created entity
	var newEntity T
	if err := tx.Table(tableName).First(&newEntity, fmt.Sprintf("%s = ?", primaryKey), tx.Statement.ReflectValue.FieldByName(primaryKey).Interface()).Error; err != nil {
		tx.Rollback()
		return nil, common.ErrDB(err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, common.ErrDB(err)
	}

	return &newEntity, nil
}

func (s *GenericStore[T]) Delete(ctx context.Context, id string) error {
	var entity T
	var tableName string

	// Check if T has TableName method
	if tableNamer, ok := any(new(T)).(interface{ TableName() string }); ok {
		tableName = tableNamer.TableName()
	} else {
		// Fallback to default table name if TableName method is not implemented
		entityType := reflect.TypeOf(new(T)).Elem().Name()
		tableName = s.db.NamingStrategy.TableName(entityType)
	}

	if err := s.db.WithContext(ctx).Table(tableName).Where("id=?", id).Delete(&entity).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}

func (s *GenericStore[T]) Update(ctx context.Context, updateData *T) (*T, error) {
	// Get the ID field value
	idValue := reflect.ValueOf(updateData).Elem().FieldByName("Id")
	if !idValue.IsValid() {
		return nil, common.ErrMissingRequiredField("Id")
	}

	// Check if T has TableName method
	var tableName string
	if tableNamer, ok := any(new(T)).(interface{ TableName() string }); ok {
		tableName = tableNamer.TableName()
	} else {
		// Fallback to default table name if TableName method is not implemented
		entityType := reflect.TypeOf(new(T)).Elem().Name()
		tableName = s.db.NamingStrategy.TableName(entityType)
	}

	// Start a transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, common.ErrDB(tx.Error)
	}

	// Update the entity
	if err := tx.Table(tableName).Where("id = ?", idValue.Interface()).Updates(updateData).Error; err != nil {
		tx.Rollback()
		return nil, common.ErrDB(err)
	}

	// Fetch the updated entity
	var updatedEntity T
	if err := tx.Table(tableName).Where("id = ?", idValue.Interface()).First(&updatedEntity).Error; err != nil {
		tx.Rollback()
		return nil, common.ErrDB(err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, common.ErrDB(err)
	}

	return &updatedEntity, nil
}

func (s *GenericStore[T]) FindWithConditions(ctx context.Context, conditions map[string]interface{}, paging *common.Paging, orderClauses []string, moreInfo ...string) ([]T, error) {
	// Determine the table name
	var tableName string
	if tableNamer, ok := any(new(T)).(interface{ TableName() string }); ok {
		tableName = tableNamer.TableName()
	} else {
		entityType := reflect.TypeOf(new(T)).Elem().Name()
		tableName = s.db.NamingStrategy.TableName(entityType)
	}

	// Prepare the query with the provided conditions
	db := s.db.Table(tableName).Where(conditions)

	// Apply optional preloading for related data
	for _, info := range moreInfo {
		db = db.Preload(info)
	}

	// Get the total count for pagination
	if err := db.Count(&paging.Total).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	// Apply pagination conditions
	if paging.FakeCusor != "" {
		if uid, err := common.FromBase58(paging.FakeCusor); err == nil {
			db = db.Where("id < ?", uid.GetLocalID())
		}
	} else {
		db = db.Offset((paging.Page - 1) * paging.Limit)
	}

	// Apply order clauses
	if len(orderClauses) > 0 {
		for _, orderClause := range orderClauses {
			db = db.Order(orderClause)
		}
	} else {
		db = db.Order("id desc")
	}

	// Execute the query with pagination and sorting
	var entities []T
	if err := db.
		Limit(paging.Limit).
		Find(&entities).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	// Set the NextCursor based on the last item in the result set
	if len(entities) > 0 {
		lastItem := entities[len(entities)-1]
		if idField := reflect.ValueOf(lastItem).FieldByName("Id"); idField.IsValid() {
			paging.NextCursor = common.NewUID(uint32(idField.Int()), int(idField.Int()), 1).String()
		}
	}

	return entities, nil
}
