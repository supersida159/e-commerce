package entities_carts

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/src/product/entities_product"
)

// Cart represents a shopping cart.

const EntityName = "Cart"

type Cart struct {
	common.SQLModel `json:",inline"`
	UserID          int              `gorm:"index;column:UserID;unique" json:"user_id"`
	Items           CartQuantityList `gorm:"column:products;type:json" json:"ProductQuantity"`
}

// CartItem represents an item in the shopping cart.

type CartQuantity struct {
	Product  entities_product.Product `json:"product"`
	Quantity int                      `json:"quantity"`
}
type CartQuantityList []*CartQuantity

func (pql CartQuantityList) Value() (driver.Value, error) {
	// Marshal the slice of pointers to JSON
	data, err := json.Marshal(pql)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Scan implements the sql.Scanner interface for []*ProductQuantity.
func (pql *CartQuantityList) Scan(value interface{}) error {
	// Check if the value is nil
	if value == nil {
		return errors.New("value is nil")
	}

	// Check if the value is a []byte
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("value is not []byte")
	}

	// Create a temporary slice to unmarshal into
	var tempSlice []*CartQuantity

	// Unmarshal the JSON data into the temporary slice
	err := json.Unmarshal(bytes, &tempSlice)
	if err != nil {
		return err
	}

	// Assign the temporary slice to the receiver
	*pql = tempSlice

	return nil
}

func (c Cart) TableName() string {
	return "Cart"
}
