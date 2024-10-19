package order_request

import (
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
)

type CreateOrderRequest struct {
	CustomerName  string                       `json:"customer_name" validate:"required"`
	CustomerPhone string                       `json:"customer_phone" validate:"required"`
	CartID        int                          `json:"cart_id" validate:"required"`
	Notes         string                       `json:"notes" validate:""`
	AddressID     int                          `json:"address_id" validate:"required"`
	Shipping      entities_orders.ShippingInfo `json:"shipping" validate:"required"`
}
