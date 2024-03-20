package subscriber

// import (
// 	"context"

// 	"github.com/supersida159/e-commerce/src/order/repository_orders"
// 	"github.com/supersida159/learningGO/component"
// 	"github.com/supersida159/learningGO/pubsub"
// 	"github.com/supersida159/learningGO/skio"
// )

// func RunDecreaseLikeCountAfterUserLikeRestaurant(appCtx component.Appcontext) consumerJob {
// 	return consumerJob{
// 		Title: "PleaceAnNewOrder",
// 		Hld: func(ctx context.Context, msg *pubsub.Message) error {
// 			store := repository_orders.NewSQLStore(appCtx.GetMainDBConnection())
// 			likedata := msg.Data().(HashRestaurantID)
// 			return store.DecreaseLikeCount(ctx, likedata.GetRestaurantID())
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
