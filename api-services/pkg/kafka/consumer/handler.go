package consumer

// import (
// 	"context"
// 	"fmt"

// 	"github.com/supersida159/e-commerce/api-services/pkg/kafka/producer"
// 	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
// 	usecase_orders "github.com/supersida159/e-commerce/api-services/src/order/usecase_order"
// )

// // OrderEventHandler implements the EventHandler interface
// type OrderEventHandler struct {
// 	orderBusiness *usecase_orders.OrderBusiness
// 	producer      *producer.OrderProducer
// }

// // NewOrderEventHandler creates a new OrderEventHandler instance
// func NewOrderEventHandler(business *usecase_orders.OrderBusiness, producer *producer.OrderProducer) *OrderEventHandler {
// 	return &OrderEventHandler{
// 		orderBusiness: business,
// 		producer:      producer,
// 	}
// }

// func (h *OrderEventHandler) HandleOrderCreated(ctx context.Context, event *entities_orders.OrderEvent) error {
// 	// Validate the order
// 	order, ok := event.Payload.(*entities_orders.Order)
// 	if !ok {
// 		return fmt.Errorf("invalid payload type for order created event")
// 	}

// 	// Create event for inventory service
// 	inventoryEvent := entities_orders.OrderEvent{
// 		Order:    *order,
// 		Payload:  order,
// 		Metadata: event.Metadata,
// 	}

// 	// Send to inventory service
// 	opts := &producer.SendMessageOptions{
// 		TargetServices: []producer.ServiceID{producer.InventoryService},
// 	}

// 	if _, err := h.producer.SendMessages(inventoryEvent, opts); err != nil {
// 		return fmt.Errorf("failed to send inventory reservation request: %w", err)
// 	}

// 	return nil
// }

// func (h *OrderEventHandler) HandleInventoryReserved(ctx context.Context, event *entities_orders.OrderEvent) error {
// 	// Update order status based on inventory reservation result
// 	success := true
// 	if event.Metadata != nil {
// 		if status, ok := event.Metadata["success"].(bool); ok {
// 			success = status
// 		}
// 	}

// 	if success {
// 		// Proceed with payment processing
// 		paymentEvent := entities_orders.OrderEvent{
// 			Order:    event.Order,
// 			Metadata: event.Metadata,
// 		}

// 		opts := &producer.SendMessageOptions{
// 			TargetServices: []producer.ServiceID{producer.SagaCentral},
// 		}

// 		if _, err := h.producer.SendMessages(paymentEvent, opts); err != nil {
// 			return fmt.Errorf("failed to send payment processing request: %w", err)
// 		}
// 	}

// 	return nil
// }

// // func (h *OrderEventHandler) HandleRollback(ctx context.Context, event *entities_orders.OrderEvent) error {
// // 	// Get the list of services that need to be rolled back
// // 	var successfulServices []producer.ServiceID
// // 	if event.Metadata != nil {
// // 		if services, ok := event.Metadata["successfulServices"].([]producer.ServiceID); ok {
// // 			successfulServices = services
// // 		}
// // 	}

// // 	// Process rollback
// // 	order, ok := event.Payload.(*entities_orders.Order)
// // 	if !ok {
// // 		return fmt.Errorf("invalid payload type for rollback event")
// // 	}

// // 	return nil
// // }
