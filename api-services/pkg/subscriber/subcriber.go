package subscriber

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/pkg/asyncjob"
	"github.com/supersida159/e-commerce/api-services/pkg/pubsub"
)

type consumerJob struct {
	Title string
	Hld   func(ctx context.Context, msg *pubsub.Message) *common.AppError
}

type consumerEngine struct {
	appCtx app_context.AppContext
	// rtEngine skio.RealTimeEngine
}

func NewEngine(appCtx app_context.AppContext,

// rtEngine skio.RealTimeEngine
) *consumerEngine {
	return &consumerEngine{
		appCtx: appCtx,
		// rtEngine: rtEngine,
	}
}

func (engine *consumerEngine) Start( /*rtEngine skio.RealtimeEngine*/ ) *common.AppError {
	// engine.startSubTopic(
	// 	pubsub.TopicUserLikeRestaurant,
	// 	false,
	// 	RunIncreaseLikeCountAfterUserLikeRestaurant(engine.appCtx),
	// )

	// engine.startSubTopic(
	// 	pubsub.TopicUserUnLikeRestaurant,
	// 	false,
	// 	RunDecreaseLikeCountAfterUserLikeRestaurant(engine.appCtx),
	// 	// EmitIncreaseLikeCountAfterUserLikeRestaurant(engine.rtEngine),
	// )
	engine.startSubExpiredKeyTopic(
		false,
		// PushNotification(engine.appCtx),
		// UpdateNotificationStatusToSent(engine.appCtx),
	)
	return nil
}

func (engine *consumerEngine) startSubTopic(topic pubsub.Topic,
	isConcurrent bool,
	consumerJobs ...consumerJob) *common.AppError {
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

func (engine *consumerEngine) startSubExpiredKeyTopic(
	isConcurrent bool,
	consumerJobs ...consumerJob) *common.AppError {

	redisClient := engine.appCtx.GetRedisClient()
	pubsubRedis := redisClient.PSubscribe(context.Background(), "__keyevent@*__:expired")
	ch := pubsubRedis.Channel()

	for _, item := range consumerJobs {
		log.Println("SetUp consumer for expired keys:", item.Title)
	}

	getJobHandler := func(job *consumerJob, msg *redis.Message, data interface{}) asyncjob.JobHandler {
		return func(ctx context.Context) *common.AppError {
			if !strings.HasPrefix(msg.Payload, "notification-") {
				return nil
			}

			message := pubsub.NewMessage(data)
			log.Printf("Running job for expired key: %s with data: %v", job.Title, data)
			return job.Hld(ctx, message)
		}
	}

	go func() {
		for msg := range ch {
			if !strings.HasPrefix(msg.Payload, "notification-") {
				continue
			}

			// Get the data once for all handlers
			preNotificationKey := "pre-" + msg.Payload
			val, err := redisClient.Get(context.Background(), preNotificationKey).Result()
			if err != nil {
				log.Printf("error getting pre-notification key: %v", err)
				continue
			}

			var dataRedis interface{}
			if err := json.Unmarshal([]byte(val), &dataRedis); err != nil {
				log.Printf("error unmarshaling data: %v", err)
				continue
			}

			// Create jobs with shared data
			jobHdlArr := make([]asyncjob.Job, len(consumerJobs))
			for i := range consumerJobs {
				jobHdl := getJobHandler(&consumerJobs[i], msg, dataRedis)
				jobHdlArr[i] = asyncjob.NewJob(jobHdl)
			}

			// Run all jobs
			group := asyncjob.NewGroup(isConcurrent, jobHdlArr...)
			if err := group.Run(context.Background()); err != nil {
				log.Println("Error in asyncjob for expired key:", err)
				continue
			}

			// Delete key only after all jobs complete successfully
			if err := redisClient.Del(context.Background(), preNotificationKey).Err(); err != nil {
				log.Printf("error deleting pre-notification key: %v", err)
			}
		}
	}()

	return nil
}
