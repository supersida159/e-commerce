package usecase_orders

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/pkg/pubsub"
	"github.com/supersida159/e-commerce/pkg/redis"
	entities_orders "github.com/supersida159/e-commerce/src/order/entities_order"
)

type CreateOrderStore interface {
	CreateOrder(ctx context.Context, data *entities_orders.Order) error
	OrderCancelled(ctx context.Context, data *entities_orders.Order) error
}

type createOrderBiz struct {
	store  CreateOrderStore
	pubsub pubsub.PubSub
	rdb    redis.RedisWRealStore
}

func NewCreateOrderBiz(store CreateOrderStore, pubsub pubsub.PubSub, rdb redis.RedisWRealStore) *createOrderBiz {
	return &createOrderBiz{
		store:  store,
		pubsub: pubsub,
		rdb:    rdb,
	}
}

func (biz *createOrderBiz) CreateOrderBiz(ctx context.Context, data *entities_orders.Order) error {

	err := biz.store.CreateOrder(ctx, data)
	if err != nil {
		return common.ErrCannotCreateEntity(entities_orders.EntityName, err)
	}

	biz.pubsub.Publish(ctx, common.TopicOrderCreated, pubsub.NewMessage(data))
	jsonData, _ := json.Marshal(data)
	biz.rdb.Client.Set(ctx, strconv.Itoa(data.ID)+"-order", "", common.ExpireOrderTime)
	biz.rdb.Client.Set(ctx, strconv.Itoa(data.ID), jsonData, 0)
	dataID := data.ID
	go func() {
		ch, close := biz.pubsub.Subscribe(ctx, common.TopicOrderExpired)
		defer close()

		for {
			ordercancel := <-ch
			// Handle received data
			// You can add your logic here to process the received data
			fmt.Println("received data on channel:", ordercancel.Data())

			if order, ok := ordercancel.Data().(*entities_orders.Order); ok {
				// Handle received order data
				// Perform action if status is not "payment success"
				if order.ID == dataID {
					biz.store.OrderCancelled(ctx, ordercancel.Data().(*entities_orders.Order))
					// Cancel goroutine if status is "payment success"
					fmt.Println("Exiting goroutine: order status is \"cancelled order\"")
					return // Perform some action with order, or ignore it
				} else {
					continue
				}

			} else {
				// Handle data with ID and status
				// Assuming ordercancel.Data() returns a struct with ID and Status fields
				order := ordercancel.Data().(struct {
					ID     int
					Status string
				})
				if order.ID == dataID {
					if order.Status == "payment success" {
						// Perform action if status is not "payment success"
						fmt.Println("Received data with ID and status on channel:", data)
						// Perform some action with data, or ignore it
					} else {
						// Cancel goroutine if status is "payment success"
						fmt.Println("Exiting goroutine: status is \"payment success\"")
						return
					}
				} else {
					continue
				}
			}

		}

	}()

	return nil
}
