package entities_orders

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/src/product/entities_product"
)

const (
	EntityName = "Order"
)

type Order struct {
	common.SQLModel `json:",inline"` // Inline embedding of common.SQLModel struct
	UserOrderID     int              `gorm:"column:user_order_id" json:"-"`

	// User-facing fields
	CustomerName  string              `gorm:"column:customer_name" json:"customer_name"`                      // Name of the customer
	CustomerPhone string              `gorm:"column:customer_phone" json:"customer_phone"`                    // Phone number of the customer
	Products      ProductQuantityList `gorm:"column:products;type:json;references:ID" json:"ProductQuantity"` // Slice of ProductQuantity structs representing the products in the order (removed unnecessary `json` tag)
	Shipping      ShippingInfo        `gorm:"embedded" json:"shipping"`                                       // Embedded ShippingInfo struct representing shipping details
	OrderTotal    float64             `gorm:"column:order_total" json:"orderTotal"`                           // Total order cost
	Notes         string              `gorm:"column:notes" json:"notes"`                                      // Notes or comments related to the order
	Address       Address             `gorm:"column:address;type:json" json:"address"`                        // Address where the order should be delivered
	// Internal field (optional)
	OrderCancelled bool `gorm:"column:order_cancelled;default:false" json:"-"` // Indicates if the order has been cancelled (consider moving this to a separate "OrderState" struct)
}

type ProductQuantity struct {
	ProductID string                   `json:"product_id" gorm:"column:id;ForeignKey:id"` // Foreign key for Product
	Product   entities_product.Product `json:"product"`
	Quantity  int                      `json:"quantity" gorm:"column:quantity"`
}
type ProductQuantityList []*ProductQuantity

func (pql ProductQuantityList) Value() (driver.Value, error) {
	// Marshal the slice of pointers to JSON
	data, err := json.Marshal(pql)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Scan implements the sql.Scanner interface for []*ProductQuantity.
func (pql *ProductQuantityList) Scan(value interface{}) error {
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
	var tempSlice []*ProductQuantity

	// Unmarshal the JSON data into the temporary slice
	err := json.Unmarshal(bytes, &tempSlice)
	if err != nil {
		return err
	}

	// Assign the temporary slice to the receiver
	*pql = tempSlice

	return nil
}

// Value implements the driver Valuer interface.
func (o *Order) GetOrderTotal() float64 {
	totalOrder := o.Shipping.Cost
	for _, product := range o.Products {
		productCost, _ := strconv.ParseFloat(product.Product.Price, 64)
		totalOrder += productCost * float64(product.Quantity)
	}
	// rounded := math.Round(totalOrder*100) / 100

	o.OrderTotal = math.Round(totalOrder*100) / 100
	return totalOrder
}

type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
	State  string `json:"state"`
}

func (a Address) Value() (driver.Value, error) {
	// Convert the Address struct to JSON string
	addrJSON, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return string(addrJSON), nil
}

// Scan implements the Scanner interface.
func (a *Address) Scan(value interface{}) error {
	// Check if the value is nil
	if value == nil {
		return nil
	}

	// Convert the value to []byte
	addrBytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal Address value: %v", value)
	}

	// Unmarshal JSON string to Address struct
	if err := json.Unmarshal(addrBytes, &a); err != nil {
		return err
	}
	return nil
}

type ShippingInfo struct {
	Method            string
	Cost              float64
	EstimatedDelivery time.Time
	// Add any additional shipping details as needed
}

func (s ShippingInfo) Value() (driver.Value, error) {
	// Convert the ShippingInfo struct to JSON string
	shippingJSON, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return string(shippingJSON), nil
}

// Scan implements the Scanner interface.
func (s *ShippingInfo) Scan(value interface{}) error {
	// Check if the value is nil
	if value == nil {
		return nil
	}

	// Convert the value to []byte
	shippingBytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal ShippingInfo value: %v", value)
	}

	// Unmarshal JSON string to ShippingInfo struct
	if err := json.Unmarshal(shippingBytes, &s); err != nil {
		return err
	}
	return nil
}

type PlaceOrderReq struct {
	CustomerName  string                `gorm:"column:customer_name" json:"customer_name"`
	CustomerPhone string                `gorm:"column:customer_phone" json:"customer_phone"`
	Products      []*ProductQuantityReq `gorm:"embedded" json:"products"`
	Notes         string                `json:"notes" gorm:"column:notes"`
	Shipping      ShippingInfo          `gorm:"embedded" json:"shipping"`
	Address       Address               `json:"address" gorm:"column:address;type:json"`
}
type ProductQuantityReq struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func (o Order) TableName() string {
	return "orders"
}

type UpdateOrder struct {
	common.SQLModel `json:",inline"`
	CustomerName    string  `gorm:"column:customer_name" json:"customer_name"`
	CustomerPhone   string  `gorm:"column:customer_phone" json:"customer_phone"`
	Notes           string  `json:"notes" gorm:"column:notes"`
	OrderCancelled  bool    `gorm:"column:order_cancelled" json:"orderCancelled"`
	Address         Address `json:"address" gorm:"column:address;type:json"`
}

func (u *UpdateOrder) TableName() string {
	return Order{}.TableName()
}

type SoftDeleteOrder struct {
	common.SQLModel `json:",inline"`

	OrderCancelled bool `gorm:"column:order_cancelled" json:"orderCancelled"`
}

func (s *SoftDeleteOrder) TableName() string {
	return Order{}.TableName()
}
func (o *Order) Mask(hideID bool) {
	if hideID {
		o.GenUID(common.DbTypeOrder)
	} else if !hideID {
		o.DeID()
	}
}

type ListOrderReq struct {
	common.SQLModel `json:",inline"` // Inline embedding of common.SQLModel struct
	OrderCancelled  bool             `json:"orderCancelled"`
	CustomerName    string           `json:"customerName"`
	CustomerPhone   string           `json:"customerPhone"`
}

func (o *ListOrderReq) Mask(hideID bool) {
	if hideID {
		o.GenUID(common.DbTypeOrder)
	} else if !hideID {
		o.DeID()
	}
}
