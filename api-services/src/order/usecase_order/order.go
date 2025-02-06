package usecase_orders

import (
	"context"
	"fmt"
	"time"

	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/pubsub"
	entities_carts "github.com/supersida159/e-commerce/api-services/src/cart/entities_cart"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
)

// OrderStatus represents the current state of the order in the saga
type OrderStatus int

const (
	OrderStatusNew OrderStatus = iota
	OrderStatusPending
	OrderStatusProcessing
	OrderStatusFailed
	OrderStatusCancelled
	OrderStatusCompleted
)

const (
	defaultTimeout = 10 * time.Second
	maxRetries     = 3
)

type OrderStore interface {
	CreateOrder(ctx context.Context, order *entities_orders.Order) *common.AppError
}
type CartStore interface {
	GetCart(ctx context.Context, cartID int, moreInfor ...string) (*entities_carts.Cart, *common.AppError)
}
type orderBiz struct {
	orderStore OrderStore
	cartStore  CartStore
	pubsub     pubsub.PubSub
}

// NewOrderBusiness creates a new instance of OrderBusiness
func NewOrderBusiness(orderStore OrderStore, cartStore CartStore, pubsub pubsub.PubSub) *orderBiz {
	return &orderBiz{
		orderStore: orderStore,
		cartStore:  cartStore,
		pubsub:     pubsub,
	}
}

// CreateOrder initiates the order creation saga
func (b *orderBiz) CreateOrder(ctx context.Context, order *entities_orders.Order) (*string, *common.AppError) {

	// Set initial order status
	order.Status = int(OrderStatusPending)
	now := time.Now()
	order.CreatedAt = &now

	cart, err := b.cartStore.GetCart(ctx, order.UserOrderID)
	if err != nil {
		return nil, common.ErrInternalServerError(err)
	}
	order.Cart = cart
	order.Products = cart.Items

	if err := b.validateOrder(order); err != nil {
		return nil, common.ErrInvalidInputData(err)
	}

	// Create and initialize the saga event
	orderEvent := b.initializeSagaEvent(order)

	publishErr := b.pubsub.Publish(ctx, pubsub.CreateOrder, pubsub.NewMessage(orderEvent))
	if publishErr != nil {
		return nil, common.ErrInternalServerError(fmt.Errorf("failed to publish order event: %w", publishErr))
	}
	return &orderEvent.SagaID, nil
}

// // HandleOrderStatusUpdate processes status updates for an order
// func (b *orderBiz) HandleOrderStatusUpdate(ctx context.Context, event *entities_orders.OrderEvent) *common.AppError {
// 	if event == nil {
// 		return common.ErrInvalidRequestParameter(fmt.Errorf("event cannot be nil"))
// 	}

// 	ctx, cancel := context.WithTimeout(ctx, b.timeout)
// 	defer cancel()

// 	// Update order status based on saga status
// 	switch event.ServiceStatus.Status {
// 	case entities_orders.ServiceSuccess:
// 		event.Order.Status = int(OrderStatusCompleted)
// 	case entities_orders.ServiceFailed:
// 		event.Order.Status = int(OrderStatusFailed)
// 	case entities_orders.ServiceCancelled:
// 		event.Order.Status = int(OrderStatusCancelled)
// 	default:
// 		event.Order.Status = int(OrderStatusProcessing)
// 	}

// 	// Handle the status update through the orchestrator
// 	if err := b.orchestrator.HandleServiceResponse(ctx, *event); err != nil {
// 		return common.ErrInternalServerError(fmt.Errorf("failed to handle status update: %w", err))
// 	}

// 	return nil
// }

// // HandleOrderCompensation processes compensation events for failed orders
// func (b *orderBiz) HandleOrderCompensation(ctx context.Context, event *entities_orders.OrderEvent) *common.AppError {
// 	if event == nil {
// 		return common.ErrInvalidRequestParameter(fmt.Errorf("event cannot be nil"))
// 	}

// 	ctx, cancel := context.WithTimeout(ctx, b.timeout)
// 	defer cancel()

// 	// Mark the order as cancelled during compensation
// 	event.Order.Status = int(OrderStatusCancelled)
// 	now := time.Now()
// 	event.Order.UpdatedAt = &now

// 	// Process the compensation through the orchestrator
// 	if err := b.orchestrator.HandleCompensation(ctx, *event); err != nil {
// 		return common.ErrInternalServerError(fmt.Errorf("failed to handle compensation: %w", err))
// 	}

// 	return nil
// }

// Helper methods

func (b *orderBiz) validateOrder(order *entities_orders.Order) error {
	if order == nil {
		return fmt.Errorf("order cannot be nil")
	}
	if order.UserOrderID <= 0 {
		return fmt.Errorf("invalid user ID")
	}
	if len(order.Cart.Items) == 0 {
		return fmt.Errorf("order must contain at least one item")
	}
	for _, item := range order.Cart.Items {
		if item.ProductID <= 0 || item.Quantity <= 0 {
			return fmt.Errorf("invalid product ID or quantity")
		}
	}
	return nil
}

func (b *orderBiz) initializeSagaEvent(order *entities_orders.Order) *entities_orders.OrderEvent {
	sagaID := fmt.Sprintf("saga_%d_%s", order.ID, time.Now().Format("20060102150405"))
	event := entities_orders.CreateSagaStartEvent(order, sagaID)

	// Initialize service status
	event.ServiceStatus = entities_orders.ServiceStatus{
		Status:      entities_orders.ServicePending,
		CurrentStep: entities_orders.EventOrderCreated,
		LastUpdated: time.Now(),
	}

	event.RetryCount = 0
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	return &event
}

// GetOrderStatus retrieves the current status of an order
func (b *orderBiz) GetOrderStatus(ctx context.Context, orderID int64) (OrderStatus, error) {
	// Implement order status retrieval logic
	// This could involve checking a database or cache
	return OrderStatusPending, nil
}

// CancelOrder initiates the cancellation process for an order
func (b *orderBiz) CancelOrder(ctx context.Context, orderID int64) *common.AppError {
	// Implement order cancellation logic
	return nil
}
