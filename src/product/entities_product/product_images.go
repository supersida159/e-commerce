package entities_product

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/supersida159/e-commerce/common"
)

type ProductImage struct {
	Color     string        `json:"color"`
	ColorCode string        `json:"colorCode"`
	Image     *common.Image `json:"image"`
}

func (j *ProductImage) TableName() string {
	return "images"
}

func (j *ProductImage) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	var img ProductImage
	if err := json.Unmarshal(bytes, &img); err != nil {
		return err
	}
	*j = img
	return nil
}
func (j *ProductImage) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}
