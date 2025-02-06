package saga

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/pkg/asyncjob"
	"github.com/supersida159/e-commerce/api-services/pkg/kafka/consumerlocal"
	kafkaconfig "github.com/supersida159/e-commerce/api-services/pkg/kafka/kafka_config"
	"github.com/supersida159/e-commerce/api-services/pkg/kafka/producers"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
)

const (
	maxRetries         = 3
	maxSagaTime        = 30 * time.Second
	updateTimeout      = 10 * time.Second
	maxConcurrentSagas = 100 // Adjust based on your needs
)

type Step string

const (
	CreateOrder Step = "CreateOrder"
	UpdateOrder Step = "UpdateOrder"
)

type SagaStep struct {
	StepName Step
	Action   func(context.Context, *entities_orders.OrderEvent) *common.AppError
	Rollback func(context.Context, *entities_orders.OrderEvent) *common.AppError
}

type Orchestrator struct {
	producer     *producers.OrderProducer
	appCtx       app_context.AppContext
	sagaSteps    []SagaStep
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
	sagaQueue    chan entities_orders.OrderEvent
	workerPool   chan struct{}
	shutdownChan chan struct{}
	OrderChannel chan entities_orders.OrderEvent
}

func NewOrchestrator(appCtx app_context.AppContext) *Orchestrator {
	producer := appCtx.GetProducer()
	if producer == nil {
		panic("producer cannot be nil")
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Orchestrator{
		producer:     producer,
		appCtx:       appCtx,
		ctx:          ctx,
		cancel:       cancel,
		sagaQueue:    make(chan entities_orders.OrderEvent, 1000), // Buffer size
		workerPool:   make(chan struct{}, maxConcurrentSagas),
		shutdownChan: make(chan struct{}),
		sagaSteps: []SagaStep{
			{
				StepName: CreateOrder,
				Action: func(ctx context.Context, event *entities_orders.OrderEvent) *common.AppError {
					createJob := asyncjob.NewGroup(false, asyncjob.NewJob(func(ctx context.Context) *common.AppError {
						if err := producer.SendCreateOrder(ctx, *event); err != nil {
							event.RetryCount++
							return err
						}
						appCtx.GetConsumer().CreateNewUpdateChannel(event.SagaID)
						return nil
					}))

					return createJob.Run(ctx)
				},
				Rollback: func(ctx context.Context, event *entities_orders.OrderEvent) *common.AppError {
					return nil // no need to rollback at this step
				},
			},
			{
				StepName: UpdateOrder,
				Action: func(ctx context.Context, event *entities_orders.OrderEvent) *common.AppError {
					appCtx.GetConsumer().SagaStates[event.SagaID].ServiceStatus.Status = entities_orders.ServiceProcessing
					// Initialize service states to Pending for all relevant services
					event.ServiceStatus.ServiceStates = map[entities_orders.ServiceName]entities_orders.ServiceState{
						entities_orders.OrderService:     {Status: entities_orders.ServicePending},
						entities_orders.CartService:      {Status: entities_orders.ServicePending},
						entities_orders.InventoryService: {Status: entities_orders.ServicePending},
					}
					appCtx.GetConsumer().SagaStates[event.SagaID].ServiceStatus.Status = entities_orders.ServiceFailed

					updateChannel, err := appCtx.GetConsumer().GetEventChannel(event.SagaID, consumerlocal.UpdateChannel)
					if err != nil {
						return common.ErrInternalServerError(err)
					}

					// Set timeout for the entire Update step
					ctxUpdate, cancel := context.WithTimeout(ctx, updateTimeout)
					defer cancel()

					for {
						select {
						case response := <-updateChannel:
							// Merge received service states into the current event
							for service, state := range response.ServiceStatus.ServiceStates {
								if _, exists := event.ServiceStatus.ServiceStates[service]; exists {
									event.ServiceStatus.ServiceStates[service] = state
								}
							}

							if event.ServiceStatus.ServiceStates[event.CurrentService].Status == entities_orders.ServiceFailed {
								return common.ErrInternalServerError(fmt.Errorf("one or more services failed during update"))
							}

							// Check if all services have completed
							allCompleted := true
							for _, state := range event.ServiceStatus.ServiceStates {
								if state.Status != entities_orders.ServiceSuccess {
									allCompleted = false
									break
								}
							}

							if allCompleted {
								return nil // Proceed to next step
							}

						case <-ctxUpdate.Done():
							// Handle timeout, check for pending services

							var pendingServices []string
							anyFailed := false
							for service, state := range event.ServiceStatus.ServiceStates {
								switch state.Status {
								case entities_orders.ServicePending:
									pendingServices = append(pendingServices, string(service))
								case entities_orders.ServiceFailed:
									anyFailed = true
								}
							}

							if len(pendingServices) > 0 {
								return common.ErrInternalServerError(
									fmt.Errorf("update timeout with pending services: %v", pendingServices),
								)
							}

							if anyFailed {
								return common.ErrInternalServerError(fmt.Errorf("update failed with service errors"))
							}

							return common.ErrInternalServerError(fmt.Errorf("update timeout")) // All succeeded before timeout

						case <-ctx.Done():

							return common.ErrInternalServerError(ctx.Err())
						}
					}
				},
				Rollback: func(ctx context.Context, event *entities_orders.OrderEvent) *common.AppError {
					var failedServices []string
					if event.ServiceStatus.ServiceStates[entities_orders.OrderService].Status == entities_orders.ServiceSuccess {
						failedServices = append(failedServices, string(kafkaconfig.KafkaTopics.RollbackOrder))
					}
					if event.ServiceStatus.ServiceStates[entities_orders.CartService].Status == entities_orders.ServiceSuccess {
						failedServices = append(failedServices, string(kafkaconfig.KafkaTopics.RollbackCart))
					}
					if event.ServiceStatus.ServiceStates[entities_orders.InventoryService].Status == entities_orders.ServiceSuccess {
						failedServices = append(failedServices, string(kafkaconfig.KafkaTopics.RollbackInventory))
					}
					if len(failedServices) == 0 {
						return nil
					}
					return producer.SendToMultipleTopics(ctx, *event, failedServices)
				},
			},
		},
	}
}

// Start initializes the orchestrator workers
func (o *Orchestrator) Start() {
	// Start order listener
	go o.listenForOrders()

	// Start worker pool
	for i := 0; i < maxConcurrentSagas; i++ {
		o.wg.Add(1)
		go o.sagaWorker()
	}
}

// Stop gracefully shuts down the orchestrator
func (o *Orchestrator) Stop() {
	close(o.shutdownChan)
	o.cancel()
	close(o.sagaQueue)
	o.wg.Wait()
}

// Improved listener with worker pool
func (o *Orchestrator) listenForOrders() {
	orderChannel := o.OrderChannel

	for {
		select {
		case event, ok := <-orderChannel:
			if !ok {
				log.Println("Order channel closed")
				return
			}
			select {
			case o.sagaQueue <- event:
			case <-o.ctx.Done():
				return
			}
		case <-o.shutdownChan:
			return
		}
	}
}

// Worker function to process sagas
func (o *Orchestrator) sagaWorker() {
	defer o.wg.Done()

	for {
		select {
		case event, ok := <-o.sagaQueue:
			if !ok {
				return
			}
			o.workerPool <- struct{}{} // Acquire worker slot
			err := o.processSaga(o.ctx, &event)
			<-o.workerPool // Release worker slot

			if err != nil {
				log.Printf("Saga processing failed: %v", err.RootErr.Error())
			}
		case <-o.shutdownChan:
			return
		}
	}
}

// Enhanced process with retry logic
func (o *Orchestrator) processSagaWithRetry(ctx context.Context, event *entities_orders.OrderEvent) error {
	for retry := 0; retry < maxRetries; retry++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := o.processSaga(ctx, event)
			if err == nil {
				return nil
			}

			if retry < maxRetries-1 {
				log.Printf("Retrying saga (attempt %d/%d)", retry+1, maxRetries)
				time.Sleep(time.Duration(retry+1) * time.Second)
			}
		}
	}
	return fmt.Errorf("saga failed after %d retries", maxRetries)
}

// Improved cleanup logic
func (o *Orchestrator) cleanupSaga(sagaID string) {
	consumersaga := o.appCtx.GetConsumer()
	consumersaga.Mu.Lock()
	defer consumersaga.Mu.Unlock()

	if _, exists := consumersaga.SagaStates[sagaID]; exists {
		delete(consumersaga.SagaStates, sagaID)
	}
}

// Modified processSaga with timeout control
func (o *Orchestrator) processSaga(ctx context.Context, event *entities_orders.OrderEvent) *common.AppError {
	sagaCtx, cancel := context.WithTimeout(ctx, maxSagaTime)
	defer cancel()

	done := make(chan struct{})
	var sagaErr *common.AppError

	go func() {
		defer close(done)
		sagaErr = o.processSagaSteps(sagaCtx, event)
	}()

	select {
	//temporary skip cleanupSaga
	case <-done:
		// o.cleanupSaga(event.SagaID)
		return sagaErr
	case <-sagaCtx.Done():
		// o.cleanupSaga(event.SagaID)
		return common.ErrInternalServerError(fmt.Errorf("saga timed out"))
	}
}

// Extracted processSagaSteps for better readability
func (o *Orchestrator) processSagaSteps(ctx context.Context, event *entities_orders.OrderEvent) *common.AppError {
	consumersaga := o.appCtx.GetConsumer()

	for _, step := range o.sagaSteps {
		// Update the current step in the saga state
		if err := event.UpdateStatus(
			entities_orders.ServiceProcessing,
			fmt.Sprintf("Processing step %s", step.StepName),
		); err != nil {
			return common.ErrInternalServerError(err)
		}

		// Update the saga state in the consumer
		consumersaga.Mu.Lock()
		consumersaga.SagaStates[event.SagaID] = event
		consumersaga.Mu.Unlock()

		// Execute the action for the current step
		if err := step.Action(ctx, event); err != nil {
			// Mark the service as failed
			if updateErr := event.UpdateStatus(
				entities_orders.ServiceFailed,
				fmt.Sprintf("Step %s failed: %v", step.StepName, err),
			); updateErr != nil {
				return common.ErrInternalServerError(updateErr)
			}

			// Update the saga state in the consumer
			consumersaga.Mu.Lock()
			consumersaga.SagaStates[event.SagaID] = event
			consumersaga.Mu.Unlock()

			// Start compensation
			if compErr := o.startCompensation(ctx, event); compErr != nil {
				return common.ErrInternalServerError(fmt.Errorf("compensation failed: %w", compErr))
			}
			return err
		}

		// Mark the service as successful
		if updateErr := event.UpdateStatus(
			entities_orders.ServiceSuccess,
			fmt.Sprintf("Step %s completed successfully", step.StepName),
		); updateErr != nil {
			return common.ErrInternalServerError(updateErr)
		}

		// Update the saga state in the consumer
		consumersaga.Mu.Lock()
		consumersaga.SagaStates[event.SagaID] = event
		consumersaga.Mu.Unlock()
	}

	// Mark the overall saga as successful
	if updateErr := event.UpdateStatus(
		entities_orders.ServiceSuccess,
		"Saga completed successfully",
	); updateErr != nil {
		return common.ErrInternalServerError(updateErr)
	}

	consumersaga.Mu.Lock()
	consumersaga.SagaStates[event.SagaID] = event
	consumersaga.Mu.Unlock()

	return nil
}

func (o *Orchestrator) startCompensation(ctx context.Context, event *entities_orders.OrderEvent) error {
	consumersaga := o.appCtx.GetConsumer()

	// Mark the saga as compensating
	if err := event.UpdateStatus(
		entities_orders.ServiceCompensating,
		"Starting compensation",
	); err != nil {
		return fmt.Errorf("failed to update service status: %w", err)
	}

	consumersaga.Mu.Lock()
	consumersaga.SagaStates[event.SagaID] = event
	consumersaga.Mu.Unlock()

	// Perform rollback for each step in reverse order
	for i := len(o.sagaSteps) - 1; i >= 0; i-- {
		step := o.sagaSteps[i]
		if err := step.Rollback(ctx, event); err != nil {
			return fmt.Errorf("compensation failed at step %s: %w", step.StepName, err)
		}

		// Update the service state after rollback
		if updateErr := event.UpdateStatus(
			entities_orders.ServiceCompensated,
			fmt.Sprintf("Compensated step %s", step.StepName),
		); updateErr != nil {
			return fmt.Errorf("failed to update service status: %w", updateErr)
		}

		consumersaga.Mu.Lock()
		consumersaga.SagaStates[event.SagaID] = event
		consumersaga.Mu.Unlock()
	}

	// Mark the saga as fully compensated
	if err := event.UpdateStatus(
		entities_orders.ServiceCompensated,
		"Saga fully compensated",
	); err != nil {
		return fmt.Errorf("failed to update service status: %w", err)
	}

	consumersaga.Mu.Lock()
	consumersaga.SagaStates[event.SagaID] = event
	consumersaga.Mu.Unlock()

	return nil
}

// Add this method to the Orchestrator struct in saga/orchestrator.go
func (o *Orchestrator) StartRollbackSingleListener(ctx context.Context) {
	consumer := o.appCtx.GetConsumer()

	for {
		select {
		case event := <-consumer.GetRollbackSingleChannel():
			o.handleRollbackSingle(ctx, event)
		case <-ctx.Done():
			return
		}
	}
}

func (o *Orchestrator) handleRollbackSingle(ctx context.Context, event entities_orders.OrderEvent) {
	consumersaga := o.appCtx.GetConsumer()

	// Get current saga state
	consumersaga.Mu.RLock()
	sagaState, exists := consumersaga.SagaStates[event.SagaID]
	consumersaga.Mu.RUnlock()
	if !exists {
		log.Printf("Saga %s not found for rollback single", event.SagaID)
		return
	}

	// Check if the service that sent the rollback had previously succeeded
	serviceName := event.CurrentService
	serviceState, ok := sagaState.ServiceStatus.ServiceStates[serviceName]
	if !ok {
		log.Printf("Service %s not found in saga %s", serviceName, event.SagaID)
		return
	}

	// Only rollback if the service was successful
	if serviceState.Status == entities_orders.ServiceSuccess {
		// Send rollback command
		if err := o.appCtx.GetProducer().SendOrderEvent(ctx, event, string(serviceName)); err != nil {
			log.Printf("Failed to send rollback to %s: %v", serviceName, err)
			return
		}

		// Update the service state to rolled back
		if err := sagaState.UpdateStatus(
			entities_orders.ServiceRollbackInitiated,
			fmt.Sprintf("Rollback initiated for %s", serviceName),
		); err != nil {
			log.Printf("Failed to update service status: %v", err)
			return
		}

		consumersaga.Mu.Lock()
		consumersaga.SagaStates[event.SagaID] = sagaState
		consumersaga.Mu.Unlock()

		log.Printf("Successfully rolled back %s for saga %s", serviceName, event.SagaID)
	} else {
		log.Printf("Skipping rollback for %s in saga %s - status was %s",
			serviceName, event.SagaID, serviceState.Status)
	}
}
func (o *Orchestrator) cleanupSagaWithDelay(sagaID string, ctx context.Context) {
	consumersaga := o.appCtx.GetConsumer()
	select {
	case <-time.After(maxSagaTime):
		consumersaga.Mu.Lock()
		if _, exists := consumersaga.SagaStates[sagaID]; exists {
			delete(consumersaga.SagaStates, sagaID)
		}
		consumersaga.Mu.Unlock()
	case <-ctx.Done():
		return
	}
}

func (o *Orchestrator) GetSagaQueue() chan entities_orders.OrderEvent {
	return o.sagaQueue
}
func (o *Orchestrator) GetOrderStatus(sagaID string) (*entities_orders.OrderEvent, error) {
	consumersaga := o.appCtx.GetConsumer()
	consumersaga.Mu.RLock()
	defer consumersaga.Mu.RUnlock()

	if event, exists := consumersaga.SagaStates[sagaID]; exists {
		return event, nil
	}
	return nil, fmt.Errorf("order %s not found", sagaID)
}
func (o *Orchestrator) SubscribeToOrderUpdates(sagaID string, conn *websocket.Conn) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			consumersaga := o.appCtx.GetConsumer()
			consumersaga.Mu.RLock()
			event, exists := consumersaga.SagaStates[sagaID]
			consumersaga.Mu.RUnlock()

			if !exists {
				conn.WriteJSON(map[string]string{"status": "not_found"})
				return
			}

			if event.ServiceStatus.Status == entities_orders.ServiceSuccess ||
				event.ServiceStatus.Status == entities_orders.ServiceFailed ||
				event.ServiceStatus.Status == entities_orders.ServiceCancelled ||
				event.ServiceStatus.Status == entities_orders.ServiceCompensated {
				conn.WriteJSON(map[string]interface{}{
					"status":  event.ServiceStatus.Status.String(),
					"message": "place order completed",
					"data":    event.Order,
				})
				return
			}
		}
	}
}
