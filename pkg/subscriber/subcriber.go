package subscriber

// import (
// 	"context"
// 	"log"

// 	"github.com/supersida159/learningGO/common"
// 	"github.com/supersida159/learningGO/component"
// 	"github.com/supersida159/learningGO/component/asyncjob"
// 	"github.com/supersida159/learningGO/pubsub"
// 	"github.com/supersida159/learningGO/skio"
// )

// type consumerJob struct {
// 	Title string
// 	Hld   func(ctx context.Context, msg *pubsub.Message) error
// }

// type consumerEngine struct {
// 	appCtx   component.Appcontext
// 	rtEngine skio.RealTimeEngine
// }

// func NewEngine(appCtx component.Appcontext, rtEngine skio.RealTimeEngine) *consumerEngine {
// 	return &consumerEngine{
// 		appCtx:   appCtx,
// 		rtEngine: rtEngine,
// 	}
// }

// func (engine *consumerEngine) Start( /*rtEngine skio.RealtimeEngine*/ ) error {
// 	engine.startSubTopic(
// 		common.TopicUserLikeRestaurant,
// 		true,
// 		RunIncreaseLikeCountAfterUserLikeRestaurant(engine.appCtx),
// 		EmitIncreaseLikeCountAfterUserLikeRestaurant(engine.rtEngine),
// 	)

// 	engine.startSubTopic(
// 		common.TopicUserUnLikeRestaurant,
// 		true,
// 		RunDecreaseLikeCountAfterUserLikeRestaurant(engine.appCtx),
// 		EmitIncreaseLikeCountAfterUserLikeRestaurant(engine.rtEngine),
// 	)
// 	return nil
// }

// func (engine *consumerEngine) startSubTopic(topic pubsub.Topic,
// 	isConcurrent bool,
// 	consumerJobs ...consumerJob) error {
// 	c, _ := engine.appCtx.GetPubSub().Subscribe(context.Background(), topic)
// 	for _, item := range consumerJobs {
// 		log.Println("SetUp comsumber for :", item.Title)
// 	}

// 	getJobHandler := func(job *consumerJob, msg *pubsub.Message) asyncjob.JobHandler {
// 		return func(ctx context.Context) error {
// 			log.Println("Running job for : ", job.Title, " with message: ", msg.Data())
// 			return job.Hld(ctx, msg)
// 		}

// 	}
// 	go func() {
// 		msg := <-c

// 		jobHdlArr := make([]asyncjob.Job, len(consumerJobs))
// 		for i := range consumerJobs {
// 			jobHdl := getJobHandler(&consumerJobs[i], msg)
// 			jobHdlArr[i] = asyncjob.NewJob(jobHdl)

// 		}

// 		group := asyncjob.NewGroup(isConcurrent, jobHdlArr...)
// 		if err := group.Run(context.Background()); err != nil {
// 			log.Println(err)
// 		}
// 	}()
// 	return nil
// }
