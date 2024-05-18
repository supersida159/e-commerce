package subscriber

import (
	"context"

	"github.com/supersida159/e-commerce/pkg/app_context"
	"github.com/supersida159/e-commerce/pkg/pubsub"
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
