package dbs

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

func NewDatabase(url string) (*Database, error) {
	database := &Database{} // Create an instance of the Database struct

	var err error
	database.db, err = gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return database, nil
}
func (d *Database) AutoMigrate(models ...interface{}) error {
	return d.db.AutoMigrate(models...)
}

func (d *Database) GetDB() *gorm.DB {
	return d.db
}
