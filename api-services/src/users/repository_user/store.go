package repository_user

import "gorm.io/gorm"

type sqlStore struct {
	db *gorm.DB
}

// FindDataByCondition implements restaurantbiz.UpdateRestaurantStore.

func NewSQLStore(db *gorm.DB) *sqlStore {
	return &sqlStore{db: db}
}

func (d *sqlStore) AutoMigrate(models ...interface{}) error {
	return d.db.AutoMigrate(models...)
}
