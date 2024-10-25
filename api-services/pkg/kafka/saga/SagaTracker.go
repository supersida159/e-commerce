package saga

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type SagaTracker struct {
	sagaRepo    port.SagaStateRepository
	orderRepo   port.OrderRepository
	msgProducer port.MessageProducer
	mu          sync.RWMutex
}

func NewSagaTracker(
	sagaRepo port.SagaStateRepository,
	orderRepo port.OrderRepository,
	msgProducer port.MessageProducer,
) *SagaTracker {
	return &SagaTracker{
		sagaRepo:    sagaRepo,
		orderRepo:   orderRepo,
		msgProducer: msgProducer,
	}
}

func (st *SagaTracker) HandleEvent(ctx context.Context, event SagaEvent) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	// Get current saga state
	state, err := st.sagaRepo.FindByID(ctx, event.SagaID)
	if err != nil {
		return fmt.Errorf("failed to find saga state: %w", err)
	}

	// Update step status
	step := state.Steps[event.ServiceName]
	step.Status = event.Status
	step.Data = event.Data
	step.Error = event.Error
	step.Timestamp = time.Now()
	state.Steps[event.ServiceName] = step

	// Check if saga is complete or failed
	state.Status = st.determineSagaStatus(state.Steps)
	state.UpdatedAt = time.Now()

	// Save updated state
	if err := st.sagaRepo.Update(ctx, state); err != nil {
		return fmt.Errorf("failed to update saga state: %w", err)
	}

	// Handle saga completion or failure
	if state.Status == SagaStatusCompleted {
		return st.handleSagaCompletion(ctx, event.SagaID)
	} else if state.Status == SagaStatusFailed {
		return st.handleSagaFailure(ctx, state)
	}

	return nil
}

func (st *SagaTracker) determineSagaStatus(steps map[string]SagaStep) SagaStatus {
	for _, step := range steps {
		if step.Status == StepStatusFailed {
			return SagaStatusFailed
		}
		if step.Status != StepStatusComplete {
			return SagaStatusInProgress
		}
	}
	return SagaStatusCompleted
}

func (st *SagaTracker) handleSagaCompletion(ctx context.Context, sagaID string) error {
	// Update order status
	order, err := st.orderRepo.FindByID(ctx, sagaID)
	if err != nil {
		return err
	}

	order.Status = OrderStatusComplete
	order.UpdatedAt = time.Now()

	return st.orderRepo.Update(ctx, order)
}

func (st *SagaTracker) handleSagaFailure(ctx context.Context, state *saga.SagaState) error {
	// Send compensation commands for completed steps
	for serviceName, step := range state.Steps {
		if step.Status == StepStatusComplete {
			compensationCmd := CompensationCommand{
				SagaID:      state.SagaID,
				ServiceName: serviceName,
				Data:        step.Data,
			}
			topic := fmt.Sprintf("%s.compensation", serviceName)
			if err := st.msgProducer.SendCompensation(ctx, topic, compensationCmd); err != nil {
				return err
			}
		}
	}

	// Update order status
	order, err := st.orderRepo.FindByID(ctx, state.SagaID)
	if err != nil {
		return err
	}

	order.Status = OrderStatusFailed
	order.UpdatedAt = time.Now()

	return st.orderRepo.Update(ctx, order)
}
