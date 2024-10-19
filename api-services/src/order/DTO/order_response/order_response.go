package order_response

import (
	"github.com/supersida159/e-commerce/api-services/common"
	entities_carts "github.com/supersida159/e-commerce/api-services/src/cart/entities_cart"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
	"github.com/supersida159/e-commerce/api-services/src/users/entities_user"
)

type OrderResponse struct {
	common.SQLModel `json:",inline"`             // Inline embedding of common.SQLModel struct
	UserOrderID     int                          ` json:"user_oder" validate:"required"`
	CustomerName    string                       `json:"customer_name" validate:"required"`  // Name of the customer
	CustomerPhone   string                       `json:"customer_phone" validate:"required"` // Phone number of the customer
	Cart            *entities_carts.Cart         `json:"cart" validate:"required"`           // Cart associated with the order
	CartID          int                          `json:"cart_id" validate:"required"`        // ID of the cart associated with the order
	Shipping        entities_orders.ShippingInfo `json:"shipping" validate:"required"`       // Embedded ShippingInfo struct representing shipping details
	OrderTotal      float64                      `json:"orderTotal" validate:"required"`     // Total order cost
	Notes           string                       `json:"notes" validate:""`                  // Notes or comments related to the order
	AddressID       int                          `json:"address_id" validate:"required"`     // ID of the address associated with the order
	Address         *entities_user.Address       `json:"address" validate:"-"`               // Address where the order should be delivered
	OrderCancelled  bool                         `json:"order_cancelled" validate:""`        // Indicates if the order has been cancelled
}
