package entities_orders

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/supersida159/e-commerce/api-services/common"
	entities_carts "github.com/supersida159/e-commerce/api-services/src/cart/entities_cart"
	"github.com/supersida159/e-commerce/api-services/src/product/entities_product"
	"github.com/supersida159/e-commerce/api-services/src/users/entities_user"
)

const (
	EntityName = "Order"
)

type Order struct {
	common.SQLModel `json:",inline"` // Inline embedding of common.SQLModel struct
	UserOrderID     int              `gorm:"column:user_order_id" json:"-"`
	BusinessID      string           `gorm:"uniqueIndex;not null" json:"business_id"` // UUID for distributed systems, this is using for kafka

	// User-facing fields
	CustomerName  string                       `gorm:"column:customer_name" json:"customer_name"`               // Name of the customer
	CustomerPhone string                       `gorm:"column:customer_phone" json:"customer_phone"`             // Phone number of the customer
	Products      []*entities_product.CartItem `gorm:"column:products;type:json;references:ID" json:"products"` // Slice of ProductQuantity structs representing the products in the order (removed unnecessary `json` tag)
	Cart          *entities_carts.Cart         `gorm:"foreignKey:CartID" json:"cart"`
	CartID        int                          `gorm:"column:cart_id" json:"cart_id"`        // ID of the cart associated with the order
	Shipping      ShippingInfo                 `gorm:"embedded" json:"shipping"`             // Embedded ShippingInfo struct representing shipping details
	OrderTotal    float64                      `gorm:"column:order_total" json:"orderTotal"` // Total order cost
	Notes         string                       `gorm:"column:notes" json:"notes"`
	AddressID     int                          `gorm:"column:address_id;default:1" json:"address_id"` // Notes or comments related to the order
	Address       *entities_user.Address       `gorm:"foreignKey:AddressID" json:"address"`           // Address where the order should be delivered
	// Internal field (optional)
	OrderCancelled bool `gorm:"column:order_cancelled;default:false" json:"order_cancelled"` // Indicates if the order has been cancelled (consider moving this to a separate "OrderState" struct)
}

// Value implements the driver Valuer interface.
func (o *Order) GetOrderTotal() float64 {
	totalOrder := o.Shipping.Cost
	for _, product := range o.Cart.Items {
		productCost, _ := strconv.ParseFloat(product.Product.Price, 64)
		totalOrder += productCost * float64(product.Quantity)
	}
	// rounded := math.Round(totalOrder*100) / 100

	o.OrderTotal = math.Round(totalOrder*100) / 100
	return totalOrder
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

// type PlaceOrderReq struct {
// 	CustomerName  string                 `gorm:"column:customer_name" json:"customer_name"`
// 	CustomerPhone string                 `gorm:"column:customer_phone" json:"customer_phone"`
// 	Cart          *entities_carts.Cart   `gorm:"foreignKey:CartID" json:"cart"`
// 	CartID        int                    `gorm:"column:cart_id" json:"cart_id"` // ID of the cart associated with the order
// 	Notes         string                 `json:"notes" gorm:"column:notes"`
// 	Shipping      ShippingInfo           `gorm:"embedded" json:"shipping"`
// 	AddressID     int                    `gorm:"column:address_id;index" json:"address_id"` // Notes or comments related to the order
// 	Address       *entities_user.Address `gorm:"foreignKey:AddressID" json:"address"`       // Address where the order should be delivered
// }

func (o Order) TableName() string {
	return "orders"
}

// type UpdateOrder struct {
// 	common.SQLModel `json:",inline"`
// 	CustomerName    string                 `gorm:"column:customer_name" json:"customer_name"`
// 	CustomerPhone   string                 `gorm:"column:customer_phone" json:"customer_phone"`
// 	Notes           string                 `json:"notes" gorm:"column:notes"`
// 	OrderCancelled  bool                   `gorm:"column:order_cancelled" json:"orderCancelled"`
// 	AddressID       int                    `gorm:"column:address_id" json:"address_id"`  // Notes or comments related to the order
// 	Shipping        ShippingInfo           `gorm:"embedded" json:"shipping" type:"json"` // Corrected struct tag
// 	Address         *entities_user.Address `gorm:"foreignKey:AddressID" json:"address"`  // Address where the order should be delivered
// }

// func (u *UpdateOrder) TableName() string {
// 	return Order{}.TableName()
// }

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

func (o *Order) GetUserOrderID() int {
	return o.UserOrderID
}
