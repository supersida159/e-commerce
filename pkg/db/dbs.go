package dbs

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	db *gorm.DB
}

func NewDatabase(url string) (*Database, error) {
	fmt.Println("URL_MYSQL11", url)
	database := &Database{} // Create an instance of the Database struct

	var err error
	fmt.Println("URL_MYSQL", url)
	database.db, err = gorm.Open(mysql.Open(url), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
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
