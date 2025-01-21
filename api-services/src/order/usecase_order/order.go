package usecase_orders

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/supersida159/e-commerce/api-services/common"
// 	"github.com/supersida159/e-commerce/api-services/pkg/kafka/producers"
// 	"github.com/supersida159/e-commerce/api-services/pkg/kafka/saga"
// 	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
// )

// // OrderStatus represents the current state of the order in the saga
// type OrderStatus int

// const (
// 	OrderStatusNew OrderStatus = iota
// 	OrderStatusPending
// 	OrderStatusProcessing
// 	OrderStatusFailed
// 	OrderStatusCancelled
// 	OrderStatusCompleted
// )

// const (
// 	defaultTimeout = 10 * time.Second
// 	maxRetries     = 3
// )

// // OrderBusiness handles order-related business logic
// type OrderBusiness struct {
// 	orderProducer *producers.OrderProducer
// 	orchestrator  *saga.Orchestrator
// 	timeout       time.Duration
// }

// // OrderBusinessConfig holds configuration for OrderBusiness
// type OrderBusinessConfig struct {
// 	Producer     *producers.OrderProducer
// 	Timeout      time.Duration
// 	Orchestrator *saga.Orchestrator
// }

// // NewOrderBusiness creates a new instance of OrderBusiness
// func NewOrderBusiness(config OrderBusinessConfig) (*OrderBusiness, error) {
// 	if config.Producer == nil {
// 		return nil, fmt.Errorf("producer cannot be nil")
// 	}
// 	if config.Orchestrator == nil {
// 		return nil, fmt.Errorf("orchestrator cannot be nil")
// 	}

// 	timeout := config.Timeout
// 	if timeout == 0 {
// 		timeout = defaultTimeout
// 	}

// 	return &OrderBusiness{
// 		orderProducer: config.Producer,
// 		orchestrator:  config.Orchestrator,
// 		timeout:       timeout,
// 	}, nil
// }

// // CreateOrder initiates the order creation saga
// func (b *OrderBusiness) CreateOrder(ctx context.Context, order *entities_orders.Order) *common.AppError {
// 	if err := b.validateOrder(order); err != nil {
// 		return common.ErrInvalidInputData(err)
// 	}

// 	ctx, cancel := context.WithTimeout(ctx, b.timeout)
// 	defer cancel()

// 	// Set initial order status
// 	order.Status = int(OrderStatusPending)
// 	now := time.Now()
// 	order.CreatedAt = &now

// 	// Create and initialize the saga event
// 	_ = b.initializeSagaEvent(order)

// 	// Start the saga process through the orchestrator
// 	if err := b.orchestrator.StartOrderSaga(ctx, order); err != nil {
// 		order.Status = int(OrderStatusFailed)
// 		return common.ErrInternalServerError(fmt.Errorf("failed to start order saga: %w", err))
// 	}

// 	return nil
// }

// // HandleOrderStatusUpdate processes status updates for an order
// func (b *OrderBusiness) HandleOrderStatusUpdate(ctx context.Context, event *entities_orders.OrderEvent) *common.AppError {
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
// func (b *OrderBusiness) HandleOrderCompensation(ctx context.Context, event *entities_orders.OrderEvent) *common.AppError {
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

// // Helper methods

// func (b *OrderBusiness) validateOrder(order *entities_orders.Order) error {
// 	if order == nil {
// 		return fmt.Errorf("order cannot be nil")
// 	}
// 	if order.UserOrderID <= 0 {
// 		return fmt.Errorf("invalid user ID")
// 	}
// 	if len(order.Cart.Items) == 0 {
// 		return fmt.Errorf("order must contain at least one item")
// 	}
// 	for _, item := range order.Cart.Items {
// 		if item.ProductID <= 0 || item.Quantity <= 0 {
// 			return fmt.Errorf("invalid product ID or quantity")
// 		}
// 	}
// 	return nil
// }

// func (b *OrderBusiness) initializeSagaEvent(order *entities_orders.Order) *entities_orders.OrderEvent {
// 	sagaID := fmt.Sprintf("saga_%d_%s", order.ID, time.Now().Format("20060102150405"))
// 	event := entities_orders.CreateSagaStartEvent(order, sagaID)

// 	// Initialize service status
// 	event.ServiceStatus = entities_orders.ServiceStatus{
// 		Status:      entities_orders.ServicePending,
// 		CurrentStep: entities_orders.EventOrderCreated,
// 		LastUpdated: time.Now(),
// 	}

// 	event.RetryCount = 0
// 	event.CreatedAt = time.Now()
// 	event.UpdatedAt = time.Now()

// 	return &event
// }

// // GetOrderStatus retrieves the current status of an order
// func (b *OrderBusiness) GetOrderStatus(ctx context.Context, orderID int64) (OrderStatus, error) {
// 	// Implement order status retrieval logic
// 	// This could involve checking a database or cache
// 	return OrderStatusPending, nil
// }

// // CancelOrder initiates the cancellation process for an order
// func (b *OrderBusiness) CancelOrder(ctx context.Context, orderID int64) *common.AppError {
// 	// Implement order cancellation logic
// 	return nil
// }
