package entities_orders

import (
	"time"
)

// ServiceStatus tracks the status of each service involved in the saga
type ServiceStatus struct {
	OrderService     int       `json:"order_service"`
	InventoryService int       `json:"inventory_service"`
	CartService      int       `json:"cart_service"`
	PaymentService   int       `json:"payment_service"`
	LastUpdated      time.Time `json:"last_updated"`
	CurrentStep      string    `json:"current_step"`
	CompensatingStep string    `json:"compensating_step,omitempty"`
}

// OrderEvent represents an event in the order saga
type OrderEvent struct {
	Order         `json:",inline"`
	ServiceStatus ServiceStatus `json:"service_status"`

	Payload        interface{}            `json:"payload"`
	EventType      string                 `json:"event_type"`
	EventTimestamp time.Time              `json:"event_timestamp"`
	SagaID         string                 `json:"saga_id"`            // Unique identifier for the saga instance
	StepNumber     int                    `json:"step_number"`        // Current step in the saga
	RetryCount     int                    `json:"retry_count"`        // Number of retries for current step
	Error          string                 `json:"error,omitempty"`    // Error message if any
	Metadata       map[string]interface{} `json:"metadata,omitempty"` // Additional metadata for the event
}

// ServiceStatusNumber represents the possible states of a service
type ServiceStatusNumber int

const (
	ServiceInit ServiceStatusNumber = iota
	ServiceSentFailed
	ServicePending
	ServiceProcessing
	ServiceCompleted
	ServiceFailed
	ServiceCompensating
	ServiceCompensated
	ServiceTimedOut
)

// Event types for the saga orchestration
const (
	// Saga Commands
	EventSagaStarted      = "SAGA_STARTED"
	EventSagaCompleted    = "SAGA_COMPLETED"
	EventSagaFailed       = "SAGA_FAILED"
	EventSagaCompensating = "SAGA_COMPENSATING"
	EventSagaCompensated  = "SAGA_COMPENSATED"

	// Order Service Events
	EventOrderCreated   = "ORDER_CREATED"
	EventOrderValidated = "ORDER_VALIDATED"
	EventOrderFailed    = "ORDER_FAILED"
	EventOrderCancelled = "ORDER_CANCELLED"

	// Inventory Service Events
	EventInventoryRequested   = "INVENTORY_REQUESTED"
	EventInventoryReserved    = "INVENTORY_RESERVED"
	EventInventoryFailed      = "INVENTORY_FAILED"
	EventInventoryCompensated = "INVENTORY_COMPENSATED"

	// Cart Service Events
	EventCartLocked      = "CART_LOCKED"
	EventCartCleared     = "CART_CLEARED"
	EventCartFailed      = "CART_FAILED"
	EventCartCompensated = "CART_COMPENSATED"

	// Payment Service Events
	EventPaymentProcessing = "PAYMENT_PROCESSING"
	EventPaymentCompleted  = "PAYMENT_COMPLETED"
	EventPaymentFailed     = "PAYMENT_FAILED"
	EventPaymentRefunded   = "PAYMENT_REFUNDED"
)

// NewOrderEvent creates an OrderEvent with initialized service statuses
func NewOrderEvent(order *Order, eventType string) OrderEvent {
	return OrderEvent{
		Order: *order,
		ServiceStatus: ServiceStatus{
			OrderService:     int(ServiceInit),
			InventoryService: int(ServiceInit),
			CartService:      int(ServiceInit),
			PaymentService:   int(ServiceInit),
			LastUpdated:      time.Now(),
		},
		EventType:      eventType,
		EventTimestamp: time.Now(),
		Metadata:       make(map[string]interface{}),
	}
}

// UpdateServiceStatus updates the status of a specific service
func (e *OrderEvent) UpdateServiceStatus(service string, status ServiceStatusNumber) {
	switch service {
	case "order":
		e.ServiceStatus.OrderService = int(status)
	case "inventory":
		e.ServiceStatus.InventoryService = int(status)
	case "cart":
		e.ServiceStatus.CartService = int(status)
	case "payment":
		e.ServiceStatus.PaymentService = int(status)
	}
	e.ServiceStatus.LastUpdated = time.Now()
}

// IsCompensating checks if any service is in compensating state
func (e *OrderEvent) IsCompensating() bool {
	return e.ServiceStatus.OrderService == int(ServiceCompensating) ||
		e.ServiceStatus.InventoryService == int(ServiceCompensating) ||
		e.ServiceStatus.CartService == int(ServiceCompensating) ||
		e.ServiceStatus.PaymentService == int(ServiceCompensating)
}

// AllServicesCompleted checks if all services have completed successfully
func (e *OrderEvent) AllServicesCompleted() bool {
	return e.ServiceStatus.OrderService == int(ServiceCompleted) &&
		e.ServiceStatus.InventoryService == int(ServiceCompleted) &&
		e.ServiceStatus.CartService == int(ServiceCompleted) &&
		e.ServiceStatus.PaymentService == int(ServiceCompleted)
}

// AddMetadata adds metadata to the event
func (e *OrderEvent) AddMetadata(key string, value interface{}) {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[key] = value
}

// Helper functions for creating specific events
func CreateSagaStartEvent(order *Order, sagaID string) OrderEvent {
	event := NewOrderEvent(order, EventSagaStarted)
	event.SagaID = sagaID
	event.StepNumber = 0
	return event
}

func CreateServiceEvent(order *Order, eventType string, sagaID string, stepNumber int) OrderEvent {
	event := NewOrderEvent(order, eventType)
	event.SagaID = sagaID
	event.StepNumber = stepNumber
	return event
}

func CreateCompensationEvent(order *Order, failedService string, sagaID string) OrderEvent {
	event := NewOrderEvent(order, EventSagaCompensating)
	event.SagaID = sagaID
	event.AddMetadata("failed_service", failedService)
	return event
}
