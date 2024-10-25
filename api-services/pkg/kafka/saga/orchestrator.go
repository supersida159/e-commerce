package saga

import (
	"context"
	"fmt"

	"github.com/supersida159/e-commerce/api-services/pkg/kafka/producers"
	"github.com/supersida159/e-commerce/api-services/pkg/redis"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
)

type OrchestratorService interface {
}

type Orchestrator struct {
	producer *producers.OrderProducer
	redis    *redis.RedisWRealStore
}

func NewOrchestrator(producer *producers.OrderProducer, redis *redis.RedisWRealStore) *Orchestrator {
	return &Orchestrator{
		producer: producer,
		redis:    redis,
	}
}

func (o *Orchestrator) HandleCreateOrderSaga(ctx context.Context, event entities_orders.OrderEvent) error {
	eventErr, err := o.producer.SendToAllServices(event)

	if err != nil {
		return err
	}

	for serviceID, err := range eventErr {
		if err != nil {
			if serviceID == producers.OrderService {
				event.ServiceStatus.OrderService = int(entities_orders.ServiceSentFailed)
			} else if serviceID == producers.InventoryService {
				event.ServiceStatus.InventoryService = int(entities_orders.ServiceSentFailed)
			} else if serviceID == producers.CartService {
				event.ServiceStatus.CartService = int(entities_orders.ServiceSentFailed)
			}

		}
	}

	err = o.redis.Set(fmt.Sprintf("order-%d", event.BusinessID), event)
	if err != nil {
		return err
	}

	return nil

}

func (o *Orchestrator) HandleUpdateOrderSaga(ctx context.Context, event entities_orders.OrderEvent) error {

	var oldEvent entities_orders.OrderEvent
	err := o.redis.Get(fmt.Sprintf("order-%d", event.BusinessID), oldEvent)
	if err != nil {
		return err
	}
	switch event.EventType {
	case "order":
		oldEvent.ServiceStatus.OrderService = event.ServiceStatus.OrderService
	case "inventory":
		oldEvent.ServiceStatus.InventoryService = event.ServiceStatus.InventoryService
	case "cart":
		oldEvent.ServiceStatus.CartService = event.ServiceStatus.CartService
	}

	err = o.redis.Set(fmt.Sprintf("order-%d", event.BusinessID), oldEvent)
	if err != nil {
		return err
	}

	return nil
}
