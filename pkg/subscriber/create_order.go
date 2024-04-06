package subscriber

// type HashRestaurantID interface {
// 	GetRestaurantID() int
// 	GetUserID() int
// }

// func IncreaseLikeCountAfterUserLikeRestaurant(appCtx component.Appcontext, ctx context.Context) error {
// 	c, _ := appCtx.GetPubSub().Subscribe(ctx, common.TopicUserLikeRestaurant)

// 	store := restaurantstore.NewSQLStore(appCtx.GetMainDBConnection())

// 	go func() {
// 		defer common.AppRecovery()
// 		for {
// 			msg := <-c
// 			likedata := msg.Data().(HashRestaurantID)
// 			_ = store.IncreaseLikeCount(ctx, likedata.GetRestaurantID())
// 		}
// 	}()
// 	return nil
// }

// func RunDeleteCartAfterOrder(appCtx component.Appcontext) consumerJob {
// 	return consumerJob{
// 		Title: "RunDeleteCartAfterOrder",
// 		Hld: func(ctx context.Context, msg *pubsub.Message) error {
// 			store := restaurantstore.NewSQLStore(appCtx.GetMainDBConnection())
// 			likedata := msg.Data().(HashRestaurantID)
// 			// fmt.Println("likedata: ", likedata)
// 			// engine.EmitToUser(likedata.GetUserID(), string(msg.Channel()), likedata)
// 			return store.IncreaseLikeCount(ctx, likedata.GetRestaurantID())
// 		},
// 	}
// }

// func EmitIncreaseLikeCountAfterUserLikeRestaurant(
// 	engine skio.RealTimeEngine) consumerJob {
// 	return consumerJob{
// 		Title: "RunDeleteCartAfterOrder",
// 		Hld: func(ctx context.Context, msg *pubsub.Message) error {
// 			likedata := msg.Data().(HashRestaurantID)
// 			return engine.EmitToUser(likedata.GetUserID(), string(msg.Channel()), likedata)
// 		},
// 	}
// }
