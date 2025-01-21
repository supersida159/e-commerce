package saga

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/pkg/asyncjob"
	"github.com/supersida159/e-commerce/api-services/pkg/kafka/consumerlocal"
	kafkaconfig "github.com/supersida159/e-commerce/api-services/pkg/kafka/kafka_config"
	"github.com/supersida159/e-commerce/api-services/pkg/kafka/producers"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
)

const (
	maxRetries     = 3
	maxWaitTime    = 30 * time.Second
	updateTimeout  = 30 * time.Millisecond //already rety in job
	redisKeyPrefix = "saga:"
)

type Step string

const (
	CreateOrder Step = "CreateOrder"
	UpdateOrder Step = "UpdateOrder"
	DeleteOrder Step = "DeleteOrder"
)

type SagaStep struct {
	StepName Step
	Action   func(context.Context, *entities_orders.OrderEvent) *common.AppError
	Rollback func(context.Context, *entities_orders.OrderEvent) *common.AppError
}

type Orchestrator struct {
	producer    *producers.OrderProducer
	appCtx      app_context.AppContext
	sagaSteps   []SagaStep
	consumer    *consumerlocal.SagaConsumer
	redisPrefix string
	mu          sync.RWMutex
}

func NewOrchestrator(producer *producers.OrderProducer, appCtx app_context.AppContext, consumer *consumerlocal.SagaConsumer) *Orchestrator {
	if producer == nil {
		panic("producer cannot be nil")
	}
	if appCtx == nil {
		panic("app context cannot be nil")
	}

	sagaSteps := []SagaStep{
		{
			StepName: CreateOrder,
			Action: func(ctx context.Context, event *entities_orders.OrderEvent) *common.AppError {
				createJob := asyncjob.NewGroup(false, asyncjob.NewJob(func(ctx context.Context) *common.AppError {
					if err := producer.SendCreateOrder(*event); err != nil {
						event.RetryCount++
						return err
					}
					return nil
				}))

				if err := createJob.Run(ctx); err != nil {
					return common.ErrInternalServerError(err)
				} else {
					sagaKey := fmt.Sprintf("saga:%s", event.SagaID)
					err := appCtx.GetCache().SetWithExpiration(sagaKey, *event, maxWaitTime)
					if err != nil {
						return common.ErrInternalServerError(err)
					}
				}
				return nil
			},
			Rollback: func(ctx context.Context, event *entities_orders.OrderEvent) *common.AppError {
				return nil
				// return producer.SendRollbackOrder(*event, []kafkaconfig.StepName{kafkaconfig.OrderServiceStep})
			},
		},
		{
			StepName: UpdateOrder,
			Action: func(ctx context.Context, event *entities_orders.OrderEvent) *common.AppError {
				updateJob := asyncjob.NewGroup(false, asyncjob.NewJob(func(ctx context.Context) *common.AppError {
					updateChannel, err := consumer.GetEventChannel(event.SagaID, consumerlocal.UpdateChannel)
					if err != nil {
						return common.ErrInternalServerError(err)
					}
					fmt.Println("updateChannel here: ", updateChannel)
					fmt.Println("event here: ", event)
					// Wait for response with timeout
					select {
					case response := <-updateChannel:
						if response.Error != "" {
							event.RetryCount++
							return common.ErrInternalServerError(fmt.Errorf(response.Error))
						}
						event.ServiceStatus = response.ServiceStatus
						return nil
					case <-time.After(updateTimeout):
						event.RetryCount++
						return common.ErrInternalServerError(fmt.Errorf("timeout"))
					case <-ctx.Done():
						return common.ErrInternalServerError(ctx.Err())
					}
				}))

				if err := updateJob.Run(ctx); err != nil {
					return common.ErrInternalServerError(err)
				} else {
					err = consumer.DeleteUpdateChannel(event.SagaID)
					if err != nil {
						return common.ErrInternalServerError(err)
					}
				}
				return nil
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
				return producer.SendToMultipleTopics(*event, failedServices)
			},
		},
	}

	return &Orchestrator{
		producer:    producer,
		appCtx:      appCtx,
		sagaSteps:   sagaSteps,
		consumer:    consumer,
		redisPrefix: redisKeyPrefix,
	}
}

func (o *Orchestrator) StartSaga(ctx context.Context, event *entities_orders.OrderEvent) error {
	// Initialize saga state
	event.ServiceStatus = entities_orders.ServiceStatus{
		Status:        entities_orders.ServiceInit,
		CurrentStep:   string(CreateOrder),
		ServiceStates: make(map[entities_orders.ServiceName]entities_orders.ServiceState),
	}

	// Save initial state
	if err := o.saveToRedis(ctx, *event); err != nil {
		return fmt.Errorf("failed to save initial state: %w", err)
	}

	return o.processSaga(ctx, event)
}

func (o *Orchestrator) processSaga(ctx context.Context, event *entities_orders.OrderEvent) error {
	for _, step := range o.sagaSteps {
		event.ServiceStatus.CurrentStep = string(step.StepName)

		// Save current step
		if err := o.UpdateRedis(ctx, *event); err != nil {
			return fmt.Errorf("failed to save step state: %w", err)
		}

		// Execute step with retry logic
		var stepErr error
		stepErr = step.Action(ctx, event)
		if stepErr == nil {
			continue
		}

		if stepErr != nil {
			event.ServiceStatus.Status = entities_orders.ServiceFailed
			if err := o.saveToRedis(ctx, *event); err != nil {
				return fmt.Errorf("failed to save failed state: %w", err)
			}
			return o.startCompensation(ctx, event)
		}
	}

	event.ServiceStatus.Status = entities_orders.ServiceSuccess
	return o.UpdateRedis(ctx, *event)
}

func (o *Orchestrator) startCompensation(ctx context.Context, event *entities_orders.OrderEvent) error {
	event.ServiceStatus.Status = entities_orders.ServiceCompensating

	// Execute rollback steps in reverse order
	for i := len(o.sagaSteps) - 1; i >= 0; i-- {
		step := o.sagaSteps[i]
		if err := step.Rollback(ctx, event); err != nil {
			return fmt.Errorf("compensation failed at step %s: %w", step.StepName, err)
		}
	}

	event.ServiceStatus.Status = entities_orders.ServiceCompensated
	return o.saveToRedis(ctx, *event)
}

func (o *Orchestrator) getCurrentState(ctx context.Context, sagaID string) (*entities_orders.OrderEvent, error) {
	key := o.getRedisKey(sagaID)
	var state entities_orders.OrderEvent
	if err := o.appCtx.GetCache().Get(key, &state); err != nil {
		return nil, fmt.Errorf("failed to get state: %w", err)
	}
	return &state, nil
}

func (o *Orchestrator) saveToRedis(ctx context.Context, state entities_orders.OrderEvent) error {
	key := o.getRedisKey(state.SagaID)
	return o.appCtx.GetCache().SetWithExpiration(key, state, maxWaitTime)
}

func (o *Orchestrator) UpdateRedis(ctx context.Context, state entities_orders.OrderEvent) error {
	key := o.getRedisKey(state.SagaID)
	return o.appCtx.GetCache().SetWithExpirationPreserve(ctx, key, state)
}

func (o *Orchestrator) getRedisKey(sagaID string) string {
	return fmt.Sprintf("%s%s", o.redisPrefix, sagaID)
}

// func (o *Orchestrator) shouldCompensate(state entities_orders.OrderEvent) bool {
// 	return time.Since(state.LastUpdated) > maxWaitTime &&
// 		state.ServiceStatus.Status != entities_orders.ServiceSuccess &&
// 		state.ServiceStatus.Status != entities_orders.ServiceCompensated
// }
