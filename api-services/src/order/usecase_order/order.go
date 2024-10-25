package usecase_orders

import (
	"context"
	"fmt"
	"time"

	"github.com/supersida159/e-commerce/api-services/pkg/kafka/producers"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
)

// OrderStatus represents the current state of the order in the saga
type OrderStatus int

const (
	OrderStatusCancelled OrderStatus = iota
	OrderStatusCreated
	OrderStatusPending
	OrderStatusFailed
	OrderStatusCompleted
)

type OrderBusiness struct {
	orderProducer *producers.OrderProducer
	timeout       time.Duration
}

func NewOrderBusiness(producer *producers.OrderProducer) *OrderBusiness {
	return &OrderBusiness{
		orderProducer: producer,
		timeout:       10 * time.Second,
	}
}

// CreateOrder initiates the order creation saga
func (b *OrderBusiness) CreateOrder(ctx context.Context, order *entities_orders.Order) error {
	ctx, cancel := context.WithTimeout(ctx, b.timeout)
	defer cancel()

	// Set initial order status
	order.Status = int(OrderStatusPending)

	// Create saga start event
	sagaID := fmt.Sprintf("saga_%d_%d", order.ID, time.Now().UnixNano())
	event := entities_orders.CreateSagaStartEvent(order, sagaID)
	event.UpdateServiceStatus("order", entities_orders.ServicePending)

	// Send to saga service first
	opts := &producers.SendMessageOptions{
		TargetServices: []producers.ServiceID{producers.SagaCentral},
	}

	_, err := b.orderProducer.SendMessages(event, opts)
	if err != nil {
		// Update order status to FAILED if we couldn't send the message
		order.Status = int(OrderStatusFailed)
		return fmt.Errorf("failed to initiate order saga: %w", err)
	}

	return nil
}

// // ProcessInventoryReservation handles the inventory reservation step
// func (b *OrderBusiness) ProcessInventoryReservation(ctx context.Context, order *entities_orders.Order, sagaID string) error {
// 	event := entities_orders.CreateServiceEvent(order, entities_orders.EventInventoryRequested, sagaID, 2)
// 	event.UpdateServiceStatus("inventory", entities_orders.ServicePending)

// 	opts := &producer.SendMessageOptions{
// 		TargetServices: []producer.ServiceID{producer.InventoryService},
// 	}

// 	results, err := b.orderProducer.SendMessages(event, opts)
// 	if err != nil {
// 		return fmt.Errorf("failed to send inventory reservation request: %w", err)
// 	}

// 	for serviceID, sendErr := range results {
// 		if sendErr != nil {
// 			return fmt.Errorf("inventory service %s failed to process request: %w", serviceID, sendErr)
// 		}
// 	}

// 	return nil
// }

// // ProcessCartUpdate handles the cart update step
// func (b *OrderBusiness) ProcessCartUpdate(ctx context.Context, order *entities_orders.Order, sagaID string) error {
// 	event := entities_orders.CreateServiceEvent(order, entities_orders.EventCartLocked, sagaID, 3)
// 	event.UpdateServiceStatus("cart", entities_orders.ServicePending)

// 	opts := &producer.SendMessageOptions{
// 		TargetServices: []producer.ServiceID{producer.CartService},
// 	}

// 	results, err := b.orderProducer.SendMessages(event, opts)
// 	if err != nil {
// 		return fmt.Errorf("failed to send cart update request: %w", err)
// 	}

// 	for serviceID, sendErr := range results {
// 		if sendErr != nil {
// 			return fmt.Errorf("cart service %s failed to process request: %w", serviceID, sendErr)
// 		}
// 	}

// 	return nil
// }

// // RollbackOrder handles saga rollback
// func (b *OrderBusiness) RollbackOrder(ctx context.Context, order *entities_orders.Order, sagaID string, failedService string) error {
// 	event := entities_orders.CreateCompensationEvent(order, failedService, sagaID)

// 	// Determine which services need compensation based on the failed service
// 	var successfulServices []producer.ServiceID
// 	switch failedService {
// 	case "inventory":
// 		successfulServices = []producer.ServiceID{producer.OrderService}
// 	case "cart":
// 		successfulServices = []producer.ServiceID{producer.OrderService, producer.InventoryService}
// 	}

// 	results, err := b.orderProducer.SendRollbackToSuccessfulServices(event, successfulServices)
// 	if err != nil {
// 		return fmt.Errorf("failed to send rollback messages: %w", err)
// 	}

// 	// Update order status to cancelled
// 	order.Status = int(OrderStatusCancelled)
// 	if err := b.orderRepo.Update(ctx, order); err != nil {
// 		return fmt.Errorf("failed to update order status during rollback: %w", err)
// 	}

// 	// Check rollback results
// 	for serviceID, sendErr := range results {
// 		if sendErr != nil {
// 			return fmt.Errorf("service %s failed to process rollback: %w", serviceID, sendErr)
// 		}
// 	}

// 	return nil
// }

// HandleServiceResponse processes responses from different services
func (b *OrderBusiness) HandleServiceResponse(ctx context.Context, orderID int, sagaID string, serviceID producers.ServiceID, success bool, event entities_orders.OrderEvent) error {
	order, err := b.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to find order: %w", err)
	}

	if !success {
		// Handle failure case
		failedService := string(serviceID)
		return b.RollbackOrder(ctx, order, sagaID, failedService)
	}

	// Handle success case
	var nextStep func(context.Context, *entities_orders.Order, string) error
	switch serviceID {
	case producers.OrderService:
		nextStep = b.ProcessInventoryReservation
	case producers.InventoryService:
		nextStep = b.ProcessCartUpdate
	case producers.CartService:
		order.Status = int(OrderStatusCompleted)
		if err := b.orderRepo.Update(ctx, order); err != nil {
			return fmt.Errorf("failed to update order status to completed: %w", err)
		}
		return nil
	}

	// Process next step if available
	if nextStep != nil {
		return nextStep(ctx, order, sagaID)
	}

	return nil
}
