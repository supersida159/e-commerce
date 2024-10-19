package goroutineinmain

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/pkg/pubsub"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
)

func RunExpireOrder(appctx *app_context.AppCtx) error {
	rd := appctx.GetCache().Client

	err := rd.ConfigSet(context.Background(), "notify-keyspace-events", "KEAx").Err()
	if err != nil {
		panic(err)
	}
	subscriber := rd.PSubscribe(context.Background(), "__keyevent@0__:expired")

	// defer pubsub.Close()

	// Go routine to receive expiration events
	go func() {
		for msg := range subscriber.Channel() {
			fmt.Println("Expired Key payload:", msg.Payload)
			res1 := strings.Split(msg.Payload, "-order")
			key := strings.Join(res1, "")
			fmt.Println("join string:", key)

			res, _ := rd.Get(context.Background(), key).Bytes()
			var order entities_orders.Order
			err = json.Unmarshal(res, &order)
			fmt.Println("res:", order)
			appctx.Pb.Publish(context.Background(), common.TopicOrderExpired, pubsub.NewMessage(&order))

		}
	}()
	return nil
}
