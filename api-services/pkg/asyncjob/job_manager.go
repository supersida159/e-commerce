package asyncjob

import (
	"context"
	"log"
	"sync"
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

func (g *group) Run(ctx context.Context) error {
	g.wg.Add(len(g.jobs))

	errChan := make(chan error, len(g.jobs))

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
	var err error

	for i := 0; i < len(g.jobs); i++ {
		if v := <-errChan; v != nil {
			err = v
		}
	}

	g.wg.Wait()

	return err
}

func (g *group) runJob(ctx context.Context, j Job) error {
	if err := j.Excute(ctx); err != nil {
		for {
			log.Println(err)
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
