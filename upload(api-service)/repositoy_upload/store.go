package repositoy_upload

import (
	"gorm.io/gorm"
)

type sqlStore struct {
	db *gorm.DB
}

// FindDataByCondition implements restaurantbiz.UpdateRestaurantStore.

func NewSQLStore(db *gorm.DB) *sqlStore {
	return &sqlStore{db: db}
}
