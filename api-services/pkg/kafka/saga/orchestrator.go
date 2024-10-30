package saga

import (
	"context"
	"fmt"
	"time"

	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/pkg/kafka/producers"
	"github.com/supersida159/e-commerce/api-services/pkg/pubsub"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
)

// SagaStep defines a single step in the saga
type SagaStep struct {
	ServiceID      producers.ServiceID
	EventType      string
	CompensateOn   []entities_orders.ServiceStatusNumber
	ForwardAction  func(context.Context, *entities_orders.OrderEvent) error
	RollbackAction func(context.Context, *entities_orders.OrderEvent) error
}

type Orchestrator struct {
	producer    *producers.OrderProducer
	appCtx      app_context.Appcontext
	sagaSteps   []SagaStep
	redisPrefix string
}

func NewOrchestrator(producer *producers.OrderProducer, appCtx app_context.Appcontext) *Orchestrator {
	orchestrator := &Orchestrator{
		producer:    producer,
		appCtx:      appCtx,
		redisPrefix: "order-saga-",
	}

	// Define the saga steps in order of execution
	orchestrator.sagaSteps = []SagaStep{
		{
			ServiceID: producers.OrderService,
			EventType: entities_orders.EventOrderCreated,
			CompensateOn: []entities_orders.ServiceStatusNumber{
				entities_orders.ServiceFailed,
				entities_orders.ServiceSentFailed,
				entities_orders.ServiceTimedOut,
			},
			ForwardAction:  orchestrator.processOrderStep,
			RollbackAction: orchestrator.rollbackOrderStep,
		},
		{
			ServiceID: producers.InventoryService,
			EventType: entities_orders.EventInventoryRequested,
			CompensateOn: []entities_orders.ServiceStatusNumber{
				entities_orders.ServiceFailed,
				entities_orders.ServiceSentFailed,
				entities_orders.ServiceTimedOut,
			},
			ForwardAction:  orchestrator.processInventoryStep,
			RollbackAction: orchestrator.rollbackInventoryStep,
		},
		{
			ServiceID: producers.CartService,
			EventType: entities_orders.EventCartLocked,
			CompensateOn: []entities_orders.ServiceStatusNumber{
				entities_orders.ServiceFailed,
				entities_orders.ServiceSentFailed,
				entities_orders.ServiceTimedOut,
			},
			ForwardAction:  orchestrator.processCartStep,
			RollbackAction: orchestrator.rollbackCartStep,
		},
	}

	return orchestrator
}

func (o *Orchestrator) StartOrderSaga(ctx context.Context, order *entities_orders.Order) error {
	// Create initial saga event
	event := entities_orders.CreateSagaStartEvent(order, generateSagaID())
	event.ServiceStatus.CurrentStep = o.sagaSteps[0].EventType

	// Send first step
	if err := o.processStep(ctx, &event, 0); err != nil {
		return fmt.Errorf("failed to start saga: %w", err)
	}

	return o.saveToRedis(ctx, event)
}

func (o *Orchestrator) HandleServiceResponse(ctx context.Context, event entities_orders.OrderEvent) error {
	var currentState entities_orders.OrderEvent
	lockKey := o.getRedisKey(event.SagaID)

	if err := o.acquireLock(ctx, lockKey); err != nil {
		return err
	}
	defer o.releaseLock(ctx, lockKey)

	if err := o.appCtx.GetCache().Get(lockKey, &currentState); err != nil {
		return fmt.Errorf("failed to get current state: %w", err)
	}

	// Update the service status based on the response
	currentState.ServiceStatus = event.ServiceStatus
	currentState.RetryCount = event.RetryCount
	currentState.Error = event.Error

	// Check if current step failed
	if o.hasStepFailed(currentState) {
		return o.startCompensation(ctx, &currentState)
	}

	// If step succeeded, move to next step
	if o.isStepComplete(currentState) {
		if currentState.StepNumber == len(o.sagaSteps)-1 {
			// All steps completed successfully
			return o.completeSaga(ctx, currentState)
		}
		// Move to next step
		currentState.StepNumber++
		currentState.ServiceStatus.CurrentStep = o.sagaSteps[currentState.StepNumber].EventType
		if err := o.processStep(ctx, &currentState, currentState.StepNumber); err != nil {
			return err
		}
	}

	return o.saveToRedis(ctx, currentState)
}

func (o *Orchestrator) HandleCompensation(ctx context.Context, event entities_orders.OrderEvent) error {
	var currentState entities_orders.OrderEvent
	lockKey := o.getRedisKey(event.SagaID)

	if err := o.acquireLock(ctx, lockKey); err != nil {
		return err
	}
	defer o.releaseLock(ctx, lockKey)

	if err := o.appCtx.GetCache().Get(lockKey, &currentState); err != nil {
		return fmt.Errorf("failed to get current state: %w", err)
	}

	// Update compensation status
	currentState.ServiceStatus = event.ServiceStatus

	// Check if compensation is complete
	if currentState.AreAllCompensationsDone() {
		return o.completeCompensation(ctx, currentState)
	}

	return o.saveToRedis(ctx, currentState)
}

// Helper methods
func (o *Orchestrator) processStep(ctx context.Context, event *entities_orders.OrderEvent, stepIndex int) error {
	step := o.sagaSteps[stepIndex]

	// Execute forward action
	if err := step.ForwardAction(ctx, event); err != nil {
		return err
	}

	// Send message to service
	opts := &producers.SendMessageOptions{
		TargetServices: []producers.ServiceID{step.ServiceID},
	}
	if _, err := o.producer.SendMessages(*event, opts); err != nil {
		event.UpdateServiceStatus(string(step.ServiceID), entities_orders.ServiceSentFailed)
		return err
	}

	event.UpdateServiceStatus(string(step.ServiceID), entities_orders.ServiceProcessing)
	return nil
}

func (o *Orchestrator) startCompensation(ctx context.Context, event *entities_orders.OrderEvent) error {
	event.EventType = entities_orders.EventSagaCompensating
	event.ServiceStatus.CompensatingStep = event.ServiceStatus.CurrentStep

	// Start compensation from the current step backwards
	for i := event.StepNumber; i >= 0; i-- {
		step := o.sagaSteps[i]
		if err := step.RollbackAction(ctx, event); err != nil {
			return fmt.Errorf("compensation failed at step %d: %w", i, err)
		}
	}

	return o.saveToRedis(ctx, *event)
}

func (o *Orchestrator) hasStepFailed(event entities_orders.OrderEvent) bool {
	currentStep := o.sagaSteps[event.StepNumber]
	status := event.ServiceStatus

	switch currentStep.ServiceID {
	case producers.OrderService:
		return o.isFailureStatus(entities_orders.ServiceStatusNumber(status.OrderService))
	case producers.InventoryService:
		return o.isFailureStatus(entities_orders.ServiceStatusNumber(status.InventoryService))
	case producers.CartService:
		return o.isFailureStatus(entities_orders.ServiceStatusNumber(status.CartService))
	default:
		return false
	}
}

func (o *Orchestrator) isFailureStatus(status entities_orders.ServiceStatusNumber) bool {
	return status == entities_orders.ServiceFailed ||
		status == entities_orders.ServiceSentFailed ||
		status == entities_orders.ServiceTimedOut
}

func (o *Orchestrator) isStepComplete(event entities_orders.OrderEvent) bool {
	currentStep := o.sagaSteps[event.StepNumber]
	status := event.ServiceStatus

	switch currentStep.ServiceID {
	case producers.OrderService:
		return status.OrderService == int(entities_orders.ServiceSuccess)
	case producers.InventoryService:
		return status.InventoryService == int(entities_orders.ServiceSuccess)
	case producers.CartService:
		return status.CartService == int(entities_orders.ServiceSuccess)
	default:
		return false
	}
}

func (o *Orchestrator) completeSaga(ctx context.Context, event entities_orders.OrderEvent) error {
	event.EventType = entities_orders.EventSagaCompleted

	// Publish success event
	o.appCtx.GetPubSub().Publish(ctx, "OrderComplete", pubsub.NewMessage(event))

	// Clean up Redis
	return o.appCtx.GetCache().Remove(o.getRedisKey(event.SagaID))
}

func (o *Orchestrator) completeCompensation(ctx context.Context, event entities_orders.OrderEvent) error {
	event.EventType = entities_orders.EventSagaCompensated

	// Publish compensation complete event
	o.appCtx.GetPubSub().Publish(ctx, "OrderCompensated", pubsub.NewMessage(event))

	// Clean up Redis
	return o.appCtx.GetCache().Remove(o.getRedisKey(event.SagaID))
}

// Step implementations
func (o *Orchestrator) processOrderStep(ctx context.Context, event *entities_orders.OrderEvent) error {
	// Implement order processing logic
	return nil
}

func (o *Orchestrator) processInventoryStep(ctx context.Context, event *entities_orders.OrderEvent) error {
	// Implement inventory processing logic
	return nil
}

func (o *Orchestrator) processCartStep(ctx context.Context, event *entities_orders.OrderEvent) error {
	// Implement cart processing logic
	return nil
}

func (o *Orchestrator) rollbackOrderStep(ctx context.Context, event *entities_orders.OrderEvent) error {
	event.UpdateServiceStatus("order", entities_orders.ServiceCompensating)
	return nil
}

func (o *Orchestrator) rollbackInventoryStep(ctx context.Context, event *entities_orders.OrderEvent) error {
	event.UpdateServiceStatus("inventory", entities_orders.ServiceCompensating)
	return nil
}

func (o *Orchestrator) rollbackCartStep(ctx context.Context, event *entities_orders.OrderEvent) error {
	event.UpdateServiceStatus("cart", entities_orders.ServiceCompensating)
	return nil
}

// Utility functions
func (o *Orchestrator) getRedisKey(sagaID string) string {
	return fmt.Sprintf("%s%s", o.redisPrefix, sagaID)
}

func (o *Orchestrator) acquireLock(ctx context.Context, key string) error {
	return o.appCtx.GetCache().LockKey(ctx, key)
}

func (o *Orchestrator) releaseLock(ctx context.Context, key string) {
	if err := o.appCtx.GetCache().UnlockKey(ctx, key); err != nil {
		fmt.Printf("failed to release lock for %s: %v\n", key, err)
	}
}

func (o *Orchestrator) saveToRedis(ctx context.Context, event entities_orders.OrderEvent) error {
	return o.appCtx.GetCache().Set(o.getRedisKey(event.SagaID), event)
}

// Helper function to generate saga ID
func generateSagaID() string {
	return fmt.Sprintf("saga-%d", time.Now().UnixNano())
}
