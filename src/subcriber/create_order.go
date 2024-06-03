package subscriber

import (
	"context"

	"github.com/supersida159/e-commerce/pkg/app_context"
	"github.com/supersida159/e-commerce/pkg/pubsub"
	"github.com/supersida159/e-commerce/pkg/skio"
	repository_carts "github.com/supersida159/e-commerce/src/cart/repository_cart"
)

type HashUserCreateOrderID interface {
	GetUserOrderID() int
}

func RunCreateNewCartAfterCreateAnOrder(appCtx app_context.Appcontext) consumerJob {
	return consumerJob{
		Title: "PleaceAnNewOrder",
		Hld: func(ctx context.Context, msg *pubsub.Message) error {
			store := repository_carts.NewSQLStore(appCtx.GetMainDBConnection())
			likedata := msg.Data().(HashUserCreateOrderID)

			return store.DeleteCart(ctx, likedata.GetUserOrderID())
		},
	}
}

func EmitCreateNewCartAfterCreateAnOrder(appCtx app_context.Appcontext, engine skio.RealTimeEngine) consumerJob {
	return consumerJob{
		Title: "Emit to user after create an order",
		Hld: func(ctx context.Context, msg *pubsub.Message) error {
			orderData := msg.Data().(HashUserCreateOrderID)
			return engine.EmitToUser(orderData.GetUserOrderID(), string(msg.Channel()), orderData)

		},
	}
}

// func RunCreateAnOrder(appCtx app_context.Appcontext) consumerJob {
// 	return consumerJob{
// 		Title: "PleaceAnNewOrder",
// 		Hld: func(ctx context.Context, msg *pubsub.Message) error {
// 			store := repository_orders.NewSQLStore(appCtx.GetMainDBConnection())

// 			return store.CreateOrder(ctx, msg.Data().(*entities_orders.Order))
// 		},
// 	}
// }

// func EmitDecreaseLikeCountAfterUserLikeRestaurant(
// 	engine skio.RealTimeEngine) consumerJob {
// 	return consumerJob{
// 		Title: "DecreaseLikeCountAfterUserLikeRestaurant",
// 		Hld: func(ctx context.Context, msg *pubsub.Message) error {
// 			likedata := msg.Data().(HashRestaurantID)
// 			return engine.EmitToUser(likedata.GetUserID(), string(msg.Channel()), likedata)
// 		},
// 	}
// }
