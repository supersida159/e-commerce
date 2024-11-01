package common

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// type Image struct {
// 	Id        int    `json:"id" gorm:"column:id;"`
// 	Url       string `json:"url" gorm:"column:url"`
// 	Width     int    `json:"width" gorm:"column:width"`
// 	Height    int    `json:"height" gorm:"column:height"`
// 	CloudName string `json:"cloud_name,omitempty" gorm:"column:clound_name"`
// 	Extension string `json:"extension,omitempty" gorm:"column:extension"`
// }

type Image struct {
	ID        int    `json:"id"`
	Url       string `json:"url"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	CloudName string `json:"cloud_name"`
	Extension string `json:"extension"`
}

func (j *Image) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	var img Image
	if err := json.Unmarshal(bytes, &img); err != nil {
		return err
	}
	*j = img
	return nil
}
func (j *Image) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

type Images []Image

func (j *Images) TableName() string {
	return "images"
}

func (j *Images) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	var img Images
	if err := json.Unmarshal(bytes, &img); err != nil {
		return err
	}
	*j = img
	return nil
}
func (j *Images) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}
