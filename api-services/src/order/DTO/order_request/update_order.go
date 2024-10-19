package order_request

import (
	"github.com/supersida159/e-commerce/api-services/common"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
)

type UpdateOrderRequest struct {
	common.SQLModel `json:",inline"`             // Inline embedding of common.SQLModel struct
	CustomerName    string                       `json:"customer_name" validate:""`
	CustomerPhone   string                       `json:"customer_phone" validate:""`
	Notes           string                       `json:"notes" validate:""`
	AddressID       int                          `json:"address_id" validate:"required"`
	Shipping        entities_orders.ShippingInfo `json:"shipping" validate:"required"`
	OrderCancelled  bool                         `json:"order_cancelled"`
}
