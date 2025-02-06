package subscriber

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/supersida159/e-commerce/create-order/common"
	"github.com/supersida159/e-commerce/create-order/pkg/app_context"
	"github.com/supersida159/e-commerce/create-order/pkg/asyncjob"
	entities "github.com/supersida159/e-commerce/create-order/src/model"
)

type OrderHandler func(ctx context.Context, event *entities.OrderEvent) *common.AppError

type Subscriber struct {
	appCtx   app_context.AppContext
	handlers map[string][]OrderHandler
	wg       sync.WaitGroup
	stopCh   chan struct{}
}

func NewSubscriber(
	appCtx app_context.AppContext,
) *Subscriber {
	return &Subscriber{
		appCtx:   appCtx,
		handlers: make(map[string][]OrderHandler),
		stopCh:   make(chan struct{}),
	}
}

func (s *Subscriber) RegisterHandler(channel string, handler OrderHandler) {
	if s.handlers[channel] == nil {
		s.handlers[channel] = make([]OrderHandler, 0)
	}
	s.handlers[channel] = append(s.handlers[channel], handler)
}

func (s *Subscriber) Start(ctx context.Context) error {
	// Start handlers for each channel type
	for channelType := range s.handlers {
		ch, err := s.appCtx.GetConsumer().GetEventChannel(channelType)
		if err != nil {
			return err
		}

		s.wg.Add(1)
		go s.handleEvents(ctx, channelType, ch)
	}

	return nil
}

func (s *Subscriber) handleEvents(ctx context.Context, channelType string, ch chan *entities.OrderEvent) {
	defer s.wg.Done()
	fmt.Println("handleEvents: ", channelType)
	for {
		select {
		case event, ok := <-ch:
			if !ok {
				return
			}

			s.processEvent(ctx, channelType, event)
		case <-ctx.Done():
			return
		case <-s.stopCh:
			return
		}
	}
}

func (s *Subscriber) processEvent(ctx context.Context, channelType string, event *entities.OrderEvent) {
	handlers := s.handlers[channelType]
	jobs := make([]asyncjob.Job, len(handlers))

	for i, handler := range handlers {
		hdlr := handler // Create new variable to avoid closure issues

		// Create job with the handler
		jobHandler := func(ctx context.Context) *common.AppError {
			return hdlr(ctx, event)
		}

		job := asyncjob.NewJob(jobHandler)

		// Customize retry durations if needed
		// job.SetRetryDuration([]time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second})

		jobs[i] = job
	}

	// Create job group - run concurrently for better performance
	group := asyncjob.NewGroup(true, jobs...)

	if err := group.Run(ctx); err != nil {
		log.Printf("Error processing event %s for channel %s: %v",
			event.SagaID, channelType, err)

		// Here you could implement compensating actions for failed jobs
		// For example, publishing to a dead letter queue
	}
}

func (s *Subscriber) Stop() {
	close(s.stopCh)
	s.wg.Wait()
}

func (s *Subscriber) GetAppContext() app_context.AppContext {
	return s.appCtx
}
