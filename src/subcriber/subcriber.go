package subscriber

import (
	"context"
	"log"

	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/pkg/app_context"
	"github.com/supersida159/e-commerce/pkg/asyncjob"
	"github.com/supersida159/e-commerce/pkg/pubsub"
	"github.com/supersida159/e-commerce/pkg/skio"
)

type consumerJob struct {
	Title string
	Hld   func(ctx context.Context, msg *pubsub.Message) error
}

type consumerEngine struct {
	appCtx   app_context.Appcontext
	rtEngine skio.RealTimeEngine
}

func NewEngine(
	appCtx app_context.Appcontext,
	rtEngine skio.RealTimeEngine) *consumerEngine {
	return &consumerEngine{
		appCtx:   appCtx,
		rtEngine: rtEngine,
	}
}

func (engine *consumerEngine) Start() error {
	engine.startSubTopic(
		common.TopicOrderCreated,
		false,
		RunCreateNewCartAfterCreateAnOrder(engine.appCtx),
		EmitCreateNewCartAfterCreateAnOrder(engine.appCtx, engine.rtEngine),
	)

	return nil
}

func (engine *consumerEngine) startSubTopic(topic pubsub.Topic,
	isConcurrent bool,
	consumerJobs ...consumerJob) error {
	c, _ := engine.appCtx.GetPubSub().Subscribe(context.Background(), topic)
	for _, item := range consumerJobs {
		log.Println("SetUp comsumber for :", item.Title)
	}

	getJobHandler := func(job *consumerJob, msg *pubsub.Message) asyncjob.JobHandler {
		return func(ctx context.Context) error {
			log.Println("Running job for : ", job.Title, " with message: ", msg.Data())
			return job.Hld(ctx, msg)
		}

	}

	go func() {
		for {
			msg := <-c

			jobHdlArr := make([]asyncjob.Job, len(consumerJobs))
			for i := range consumerJobs {
				jobHdl := getJobHandler(&consumerJobs[i], msg)
				jobHdlArr[i] = asyncjob.NewJob(jobHdl)

			}

			group := asyncjob.NewGroup(isConcurrent, jobHdlArr...)
			if err := group.Run(context.Background()); err != nil {
				log.Println(err)
			}
		}
	}()

	return nil
}
