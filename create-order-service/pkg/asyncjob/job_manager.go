package asyncjob

import (
	"context"
	"fmt"
	"sync"

	"github.com/supersida159/e-commerce/create-order/common"
)

type group struct {
	isConcurrent bool
	jobs         []Job
	wg           *sync.WaitGroup
}

func NewGroup(isConcurrent bool, jobs ...Job) *group {
	return &group{
		isConcurrent: isConcurrent,
		jobs:         jobs,
		wg:           new(sync.WaitGroup),
	}
}

func (g *group) Run(ctx context.Context) *common.AppError {
	g.wg.Add(len(g.jobs))

	errChan := make(chan *common.AppError, len(g.jobs))

	for i, _ := range g.jobs {
		if g.isConcurrent {
			go func(j Job) {
				errChan <- g.runJob(ctx, j)
				g.wg.Done()
			}(g.jobs[i])

			continue
		}
		job := g.jobs[i]
		err := g.runJob(ctx, job)
		errChan <- err
		g.wg.Done()
		if err != nil {
			break
		}
	}
	var err *common.AppError

	for i := 0; i < len(g.jobs); i++ {
		if v := <-errChan; v != nil {
			err = v
		}
	}

	g.wg.Wait()

	return err
}

func (g *group) runJob(ctx context.Context, j Job) *common.AppError {
	if err := j.Excute(ctx); err != nil {
		for {
			fmt.Println("err at run job:", err.RootErr)
			if j.State() == StateRetryFailed {
				return err
			}
			if err = j.Retry(ctx); err == nil {
				return nil
			}
		}
	}
	return nil
}
