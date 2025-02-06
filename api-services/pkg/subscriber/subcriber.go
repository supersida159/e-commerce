package subscriber

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/pkg/asyncjob"
	"github.com/supersida159/e-commerce/api-services/pkg/kafka/saga"
	"github.com/supersida159/e-commerce/api-services/pkg/pubsub"
	entities_orders "github.com/supersida159/e-commerce/api-services/src/order/entities_order"
)

type consumerJob struct {
	Title string
	Hld   func(ctx context.Context, msg *pubsub.Message) *common.AppError
}

type consumerEngine struct {
	appCtx       app_context.AppContext
	orchestrator *saga.Orchestrator // Add Orchestrator to the engine
}

// NewEngine initializes the consumer engine with the orchestrator.
func NewEngine(appCtx app_context.AppContext, sagaOrchestration *saga.Orchestrator) *consumerEngine {
	return &consumerEngine{
		appCtx:       appCtx,
		orchestrator: sagaOrchestration,
	}
}

// Start initializes all subscribers and starts the orchestrator.
func (engine *consumerEngine) Start() *common.AppError {
	// // Map Redis prefixes to topics
	// mapPrefixWithTopic := map[string]pubsub.Topic{
	// 	string(pubsub.OrderExpire) + ":" + string(pubsub.Saga): pubsub.UpdateOrdeExpire,
	// }

	// Listen for Redis events
	// engine.ListenRedisEventKey(mapPrefixWithTopic)

	// Subscribe to CreateOrder topic
	engine.startSubTopic(
		pubsub.CreateOrder,
		false,
		consumerJob{
			Title: "HandleCreateOrder",
			Hld:   engine.handleCreateOrder,
		},
	)

	// Start rollback listener for the orchestrator
	go engine.orchestrator.StartRollbackSingleListener(context.Background())

	// Start the orchestrator in the background
	engine.orchestrator.Start()

	return nil
}

// handleCreateOrder processes incoming CreateOrder messages and triggers the orchestrator.
func (engine *consumerEngine) handleCreateOrder(ctx context.Context, msg *pubsub.Message) *common.AppError {
	// Log the type and value of the message data
	log.Printf("Message data type: %T, value: %v", msg.Data(), msg.Data())

	var orderEvent entities_orders.OrderEvent

	switch data := msg.Data().(type) {
	case []byte:
		// Deserialize from []byte
		if err := json.Unmarshal(data, &orderEvent); err != nil {
			log.Printf("Failed to unmarshal CreateOrder message: %v", err)
			return common.ErrInvalidInputData(err)
		}
	case *entities_orders.OrderEvent:
		// Use the deserialized object directly
		orderEvent = *data
	default:
		log.Printf("Unsupported message data type: %T", msg.Data())
		return common.ErrInvalidInputData(fmt.Errorf("unsupported message data type"))
	}

	log.Printf("Received CreateOrder event: %+v", orderEvent)

	// Pass the event to the orchestrator's queue
	engine.orchestrator.GetSagaQueue() <- orderEvent
	return nil
}

// startSubTopic subscribes to a topic and processes messages using the provided handlers.
func (engine *consumerEngine) startSubTopic(topic pubsub.Topic, isConcurrent bool, consumerJobs ...consumerJob) *common.AppError {
	c, _ := engine.appCtx.GetPubSub().Subscribe(context.Background(), topic)
	for _, item := range consumerJobs {
		log.Println("SetUp consumer for:", item.Title)
	}

	getJobHandler := func(job *consumerJob, msg *pubsub.Message) asyncjob.JobHandler {
		return func(ctx context.Context) *common.AppError {
			log.Println("Running job for:", job.Title, "with message:", msg.Data())
			return job.Hld(ctx, msg)
		}
	}

	go func() {
		for msg := range c {
			jobHdlArr := make([]asyncjob.Job, len(consumerJobs))
			for i := range consumerJobs {
				jobHdl := getJobHandler(&consumerJobs[i], msg)
				jobHdlArr[i] = asyncjob.NewJob(jobHdl)
			}
			group := asyncjob.NewGroup(isConcurrent, jobHdlArr...)
			if err := group.Run(context.Background()); err != nil {
				log.Println("Error in asyncjob:", err)
			}
		}
	}()

	return nil
}

// // ListenRedisEventKey listens for Redis key events and publishes them to the appropriate topic.
// func (engine *consumerEngine) ListenRedisEventKey(mapPrefixWithTopic map[string]pubsub.Topic) {
// 	for prefixwithkeyevent, topic := range mapPrefixWithTopic {
// 		redisClient := engine.appCtx.GetRedisClient()
// 		pubsubRedis := redisClient.PSubscribe(context.Background(), prefixwithkeyevent)
// 		ch := pubsubRedis.Channel()

// 		go func(prefix string, ch <-chan *redis.Message, topic pubsub.Topic) {
// 			for msg := range ch {
// 				// Convert redis.Message to pubsub.Message
// 				pubsubMsg := pubsub.NewMessage(msg.Payload)
// 				err := engine.appCtx.GetPubSub().Publish(context.Background(), topic, pubsubMsg)
// 				if err != nil {
// 					log.Printf("Publish error for prefix %s: %v", prefix, err)
// 				}
// 			}
// 		}(prefixwithkeyevent, ch, topic)
// 	}
// }
