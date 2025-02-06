package asyncjob

import (
	"context"
	"time"

	"github.com/supersida159/e-commerce/create-order/common"
)

//job requirement:
//1. job can do something(handler)
//2. job can retry
//2.1 can config times and duration
//3. Should be stateful
//4. we should have job manager to manage job

type Job interface {
	Excute(ctx context.Context) *common.AppError
	Retry(ctx context.Context) *common.AppError
	State() JobState
	RetryIndex() int
	SetRetryDuration(time []time.Duration)
}

const (
	defaultMaxTimeOut    = 3 * time.Second
	defaultMaxRetryCount = 3
)

var defaultRetryTime = []time.Duration{1 * time.Second, 5 * time.Second, 10 * time.Second}

func getRetryTimes(n int) []time.Duration {
	if n <= 0 {
		return []time.Duration{} // Handle non-positive n
	}

	result := make([]time.Duration, 0, n)
	defaultLen := len(defaultRetryTime)

	if defaultLen > 0 {
		// Take elements from defaultRetryTime first
		take := min(n, defaultLen)
		result = append(result, defaultRetryTime[:take]...)
		remaining := n - take

		// Double the last duration for remaining elements
		if remaining > 0 {
			lastDuration := defaultRetryTime[take-1]
			for i := 0; i < remaining; i++ {
				lastDuration *= 2
				result = append(result, lastDuration)
			}
		}
	} else {
		// Edge case: If defaultRetryTime is empty, start with 1s and double
		lastDuration := 1 * time.Second
		result = append(result, lastDuration)
		for i := 1; i < n; i++ {
			lastDuration *= 2
			result = append(result, lastDuration)
		}
	}

	return result
}

type JobState int

type JobHandler func(ctx context.Context) *common.AppError

const (

	// inum auto increate
	StateInit    JobState = iota
	StateRunning          //=1
	StateFailed
	StateTimeOut
	StateCompleted
	StateRetryFailed
)

type JobConfig struct {
	MaxTimeOut time.Duration
	Retries    []time.Duration
}

func (j *JobState) String() string {

	// convert const JobState to string
	return []string{"Init", "Running", "Failed", "TimeOut", "Completed", "RetryFailed"}[*j]
}

type job struct {
	config     JobConfig
	handler    JobHandler
	state      JobState
	retryIndex int
	stopChan   chan bool
}

func NewJob(handler JobHandler, retry int) *job {
	j := job{
		config: JobConfig{
			MaxTimeOut: defaultMaxTimeOut,
			Retries:    getRetryTimes(retry),
		},
		handler:    handler,
		state:      StateInit,
		retryIndex: -1,
		stopChan:   make(chan bool),
	}

	return &j
}

func (j *job) Excute(ctx context.Context) *common.AppError {
	j.state = StateRunning

	err := j.handler(ctx)
	if err != nil {
		j.state = StateFailed
		return err
	}
	j.state = StateCompleted
	return nil
}

func (j *job) Retry(ctx context.Context) *common.AppError {
	j.retryIndex++

	time.Sleep(j.config.Retries[j.retryIndex])
	err := j.Excute(ctx)
	if err == nil {
		j.state = StateCompleted
		return nil
	}

	if j.retryIndex == len(j.config.Retries)-1 {
		j.state = StateRetryFailed
		return err
	}
	j.state = StateFailed
	return err

}

func (j *job) State() JobState {
	return j.state
}

func (j *job) RetryIndex() int {
	return j.retryIndex

}

func (j *job) SetRetryDuration(times []time.Duration) {
	j.config.Retries = times

	if len(j.config.Retries) == 0 {
		j.config.Retries = defaultRetryTime
	}
}
