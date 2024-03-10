package subscriber

// import (
// 	"context"

// 	"github.com/supersida159/learningGO/common"
// 	"github.com/supersida159/learningGO/component"
// 	"github.com/supersida159/learningGO/modules/restaurant/restaurantstore"
// 	"github.com/supersida159/learningGO/pubsub"
// 	"github.com/supersida159/learningGO/skio"
// )

// func DecreaseLikeCountAfterUserLikeRestaurant(appCtx component.Appcontext, ctx context.Context) error {
// 	c, _ := appCtx.GetPubSub().Subscribe(ctx, common.TopicUserLikeRestaurant)

// 	store := restaurantstore.NewSQLStore(appCtx.GetMainDBConnection())

// 	go func() {
// 		defer common.AppRecovery()
// 		for {
// 			msg := <-c
// 			likedata := msg.Data().(HashRestaurantID)
// 			_ = store.DecreaseLikeCount(ctx, likedata.GetRestaurantID())
// 		}
// 	}()
// 	return nil
// }

// func RunDecreaseLikeCountAfterUserLikeRestaurant(appCtx component.Appcontext) consumerJob {
// 	return consumerJob{
// 		Title: "DecreaseLikeCountAfterUserLikeRestaurant",
// 		Hld: func(ctx context.Context, msg *pubsub.Message) error {
// 			store := restaurantstore.NewSQLStore(appCtx.GetMainDBConnection())
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
