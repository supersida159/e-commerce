package entities_orders

import (
	"strconv"
	"time"

	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/src/product/entities_product"
)

const (
	EntityName = "Order"
)

type Order struct {
	common.SQLModel `json:",inline"`
	UserOrderID     int               `gorm:"column:user_order_id" json:"-"`
	CustomerName    string            `gorm:"column:customer_name" json:"customer"`
	CustomerPhone   string            `gorm:"column:customer_phone" json:"phone"`
	Products        []ProductQuantity `gorm:"embedded" json:"products"`
	// Payment         PaymentInfo                `gorm:"embedded" json:"payment"`
	Shipping       ShippingInfo `gorm:"embedded" json:"shipping"`
	OrderTotal     float64      `json:"-" gorm:"column:order_total"`
	Notes          string       `json:"notes" gorm:"column:notes"`
	Address        Address      `json:"address" gorm:"column:address;type:json"`
	OrderCancelled bool         `gorm:"column:order_cancelled;default:false" json:"orderCancelled"`
}
type ProductQuantity struct {
	Product  entities_product.Product
	Quantity int
}

func (o *Order) GetOrderTotal() float64 {
	totalOrder := o.Shipping.Cost
	for _, product := range o.Products {
		productCost, _ := strconv.ParseFloat(product.Price, 64)
		totalOrder += productCost
	}
	return totalOrder
}

type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
	State  string `json:"state"`
}
type ShippingInfo struct {
	Method            string
	Cost              float64
	EstimatedDelivery time.Time
	// Add any additional shipping details as needed
}

type PlaceOrderReq struct {
	CustomerName  string                     `gorm:"column:customer_name" json:"customer"`
	CustomerPhone string                     `gorm:"column:customer_phone" json:"phone"`
	Products      []entities_product.Product `gorm:"foreignKey:OrderID" json:"products"`
	Notes         string                     `json:"notes" gorm:"column:notes"`
	Shipping      ShippingInfo               `gorm:"embedded" json:"shipping"`
	Address       Address                    `json:"address" gorm:"column:address;type:json"`
}

func (o Order) TableName() string {
	return "orders"
}

type UpdateOrder struct {
	common.SQLModel `json:",inline"`
	CustomerName    string  `gorm:"column:customer_name" json:"customer"`
	CustomerPhone   string  `gorm:"column:customer_phone" json:"phone"`
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
