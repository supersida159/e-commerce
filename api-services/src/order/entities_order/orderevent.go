package entities_orders

import (
	"fmt"
	"time"
)

// ServiceName type for service identification
type ServiceName string

const (
	OrderService     ServiceName = "order"
	InventoryService ServiceName = "inventory"
	CartService      ServiceName = "cart"
)

type EventTypeString string

const (
	EventOrderCreate EventTypeString = "Order Create"
	EventSuccess     EventTypeString = "Order Successes"
	EventFailed      EventTypeString = "Order Failed"
	EventRollback    EventTypeString = "Order Rollback"
)

// ServiceStatus tracks the status of each service involved in the saga
type ServiceStatus struct {
	Status           ServiceStatusNumber          `json:"status"`
	CurrentStep      string                       `json:"current_step"`
	CompensatingStep string                       `json:"compensating_step,omitempty"`
	LastUpdated      time.Time                    `json:"last_updated"`
	ServiceStates    map[ServiceName]ServiceState `json:"service_states"`
}

// ServiceState represents the state of an individual service
type ServiceState struct {
	Status    ServiceStatusNumber `json:"status"`
	UpdatedAt time.Time           `json:"updated_at"`
	Error     string              `json:"error,omitempty"`
}

// OrderEvent represents an event in the order saga
type OrderEvent struct {
	Order         `json:",inline"`
	ServiceStatus ServiceStatus `json:"service_status"`

	CurrentService ServiceName            `json:"current_service"`
	EventType      string                 `json:"event_type"`
	EventTimestamp time.Time              `json:"event_timestamp"`
	SagaID         string                 `json:"saga_id"`
	StepNumber     int                    `json:"step_number"`
	RetryCount     int                    `json:"retry_count"`
	Error          string                 `json:"error,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	Version        int                    `json:"version"` // Optimistic locking
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// ServiceStatusNumber represents the possible states of a service
type ServiceStatusNumber int

const (
	ServiceInit ServiceStatusNumber = iota
	ServicePending
	ServiceProcessing
	ServiceSuccess
	ServiceFailed
	ServiceCancelled
	ServiceTimedOut
	ServiceSentFailed
	ServiceCompensating
	ServiceCompensated
	ServiceRollbackInitiated
	ServiceRollbackInProgress
	ServiceRollbackSuccess
	ServiceRollbackFailed
)

// String returns the string representation of the ServiceStatusNumber.
func (s ServiceStatusNumber) String() string {
	switch s {
	case ServiceInit:
		return "ServiceInit"
	case ServicePending:
		return "ServicePending"
	case ServiceProcessing:
		return "ServiceProcessing"
	case ServiceSuccess:
		return "ServiceSuccess"
	case ServiceFailed:
		return "ServiceFailed"
	case ServiceCancelled:
		return "ServiceCancelled"
	case ServiceTimedOut:
		return "ServiceTimedOut"
	case ServiceSentFailed:
		return "ServiceSentFailed"
	case ServiceCompensating:
		return "ServiceCompensating"
	case ServiceCompensated:
		return "ServiceCompensated"
	case ServiceRollbackInitiated:
		return "ServiceRollbackInitiated"
	case ServiceRollbackInProgress:
		return "ServiceRollbackInProgress"
	case ServiceRollbackSuccess:
		return "ServiceRollbackSuccess"
	case ServiceRollbackFailed:
		return "ServiceRollbackFailed"
	default:
		return fmt.Sprintf("Unknown ServiceStatusNumber (%d)", s)
	}
}

// Event types for the saga orchestration
const (
	// Saga Lifecycle Events
	EventSagaStarted      = "SAGA_STARTED"
	EventSagaCompleted    = "SAGA_COMPLETED"
	EventSagaFailed       = "SAGA_FAILED"
	EventSagaCompensating = "SAGA_COMPENSATING"
	EventSagaCompensated  = "SAGA_COMPENSATED"

	// Order Service Events
	EventOrderCreated     = "ORDER_CREATED"
	EventOrderValidated   = "ORDER_VALIDATED"
	EventOrderFailed      = "ORDER_FAILED"
	EventOrderCancelled   = "ORDER_CANCELLED"
	EventOrderCompensated = "ORDER_COMPENSATED"

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
	event := OrderEvent{
		Order:         *order,
		EventType:     eventType,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		ServiceStatus: newServiceStatus(),
		Metadata:      make(map[string]interface{}),
	}
	return event
}

// newServiceStatus initializes a new ServiceStatus with default values
func newServiceStatus() ServiceStatus {
	return ServiceStatus{
		Status:      ServiceInit,
		LastUpdated: time.Now(),
		ServiceStates: map[ServiceName]ServiceState{
			OrderService:     {Status: ServiceInit, UpdatedAt: time.Now()},
			InventoryService: {Status: ServiceInit, UpdatedAt: time.Now()},
			CartService:      {Status: ServiceInit, UpdatedAt: time.Now()},
		},
	}
}

// UpdateStatus updates the status of a specific service
func (e *OrderEvent) UpdateStatus(status ServiceStatusNumber, errMsg string) error {

	e.ServiceStatus.Status = status
	e.UpdatedAt = time.Now()
	e.Version++

	return nil
}

func (e *OrderEvent) UpdateServiceStatus(service ServiceName, status ServiceStatusNumber, errMsg string) error {
	state, exists := e.ServiceStatus.ServiceStates[service]
	if !exists {
		return fmt.Errorf("invalid service name: %s", service)
	}

	state.Status = status
	state.UpdatedAt = time.Now()
	state.Error = errMsg
	e.ServiceStatus.ServiceStates[service] = state
	e.ServiceStatus.LastUpdated = time.Now()
	e.UpdatedAt = time.Now()
	e.Version++

	return nil
}

// GetServiceState returns the current state of a specific service
func (e *OrderEvent) GetServiceState(service ServiceName) (ServiceState, error) {
	state, exists := e.ServiceStatus.ServiceStates[service]
	if !exists {
		return ServiceState{}, fmt.Errorf("invalid service name: %s", service)
	}
	return state, nil
}

// IsCompensating checks if any service is in compensating state
func (e *OrderEvent) IsCompensating() bool {
	for _, state := range e.ServiceStatus.ServiceStates {
		if state.Status == ServiceCompensating {
			return true
		}
	}
	return false
}

// AllServicesCompleted checks if all services have completed successfully
func (e *OrderEvent) AllServicesCompleted() bool {
	for _, state := range e.ServiceStatus.ServiceStates {
		if state.Status != ServiceSuccess {
			return false
		}
	}
	return true
}

// AreAllServicesDone checks if all services have reached a terminal state
func (e *OrderEvent) AreAllServicesDone() (bool, bool) {
	allDone := true
	allSuccess := true

	for _, state := range e.ServiceStatus.ServiceStates {
		if state.Status != ServiceSuccess && state.Status != ServiceFailed {
			allDone = false
			break
		}
		if state.Status == ServiceFailed {
			allSuccess = false
		}
	}

	return allDone, allSuccess
}

// AreAllCompensationsDone checks if all compensations have completed
func (e *OrderEvent) AreAllCompensationsDone() bool {
	for _, state := range e.ServiceStatus.ServiceStates {
		if state.Status != ServiceRollbackSuccess && state.Status != ServiceRollbackFailed {
			return false
		}
	}
	return true
}

// Helper functions for event creation
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

func CreateCompensationEvent(order *Order, failedService ServiceName, sagaID string) OrderEvent {
	event := NewOrderEvent(order, EventSagaCompensating)
	event.SagaID = sagaID
	event.AddMetadata("failed_service", string(failedService))
	return event
}

// AddMetadata adds metadata to the event
func (e *OrderEvent) AddMetadata(key string, value interface{}) {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[key] = value
}

// Validate validates the event
func (e *OrderEvent) Validate() error {
	if e.SagaID == "" {
		return fmt.Errorf("saga ID is required")
	}
	if e.EventType == "" {
		return fmt.Errorf("event type is required")
	}
	if e.Order.ID <= 0 {
		return fmt.Errorf("invalid order ID")
	}
	return nil
}
