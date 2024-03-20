package main

// import (
// 	"context"
// 	"log"
// 	"time"

// 	"github.com/supersida159/e-commerce/pkg/pubsub"
// 	"github.com/supersida159/e-commerce/pkg/pubsub/pubsublocal"
// )

// func main() {
// 	var localPB pubsub.PubSub = pubsublocal.NewPubSub()
// 	var topic pubsub.Topic = "Order Created"

// 	sub1, close1 := localPB.Subscribe(context.Background(), topic)
// 	sub2, _ := localPB.Subscribe(context.Background(), topic)
// 	localPB.Publish(context.Background(), topic, pubsub.NewMessage("1"))
// 	localPB.Publish(context.Background(), topic, pubsub.NewMessage("2"))
// 	go func() {
// 		for {
// 			log.Println("con1", (<-sub1).Data())
// 			time.Sleep(time.Millisecond * 400)
// 		}
// 	}()
// 	go func() {
// 		for {
// 			log.Println("con2", (<-sub2).Data())
// 			time.Sleep(time.Millisecond * 400)
// 		}
// 	}()
// 	time.Sleep(time.Second * 3)
// 	close1()
// 	localPB.Publish(context.Background(), topic, pubsub.NewMessage("3"))
// 	time.Sleep(time.Second * 2)
// }
