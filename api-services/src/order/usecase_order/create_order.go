package usecase_orders

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/pubsub"
	"github.com/supersida159/e-commerce/api-services/pkg/redis"
	usecase_carts "github.com/supersida159/e-commerce/api-services/src/cart/usecase_cart"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
)

type CreateOrderStore interface {
	CreateOrder(ctx context.Context, data *entities_orders.Order) *common.AppError
	// OrderCancelled(ctx context.Context, Id string) *common.AppError
}

type createOrderBiz struct {
	storeOrder CreateOrderStore
	storeCart  usecase_carts.CreateCartStore
	pubsub     pubsub.PubSub
	rdb        redis.RedisWRealStore
}

func NewCreateOrderBiz(storeOrder CreateOrderStore, storeCart usecase_carts.CreateCartStore, pubsub pubsub.PubSub, rdb redis.RedisWRealStore) *createOrderBiz {
	return &createOrderBiz{
		storeOrder: storeOrder,
		storeCart:  storeCart,
		pubsub:     pubsub,
		rdb:        rdb,
	}
}

func (biz *createOrderBiz) CreateOrderBiz(ctx context.Context, data *entities_orders.Order) *common.AppError {

	err := biz.storeOrder.CreateOrder(ctx, data)
	if err != nil {
		return common.ErrDB(err)
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
					biz.storeOrder.OrderCancelled(ctx, ordercancel.Data().(*entities_orders.Order))
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
