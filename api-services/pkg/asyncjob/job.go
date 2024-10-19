package asyncjob

import (
	"context"
	"time"
)

//job requirement:
//1. job can do something(handler)
//2. job can retry
//2.1 can config times and duration
//3. Should be stateful
//4. we should have job manager to manage job

type Job interface {
	Excute(ctx context.Context) error
	Retry(ctx context.Context) error
	State() JobState
	RetryIndex() int
	SetRetryDuration(time []time.Duration)
}

const (
	defaultMaxTimeOut    = 3 * time.Second
	defaultMaxRetryCount = 3
)

var defaultRetryTime = []time.Duration{1 * time.Second, 5 * time.Second, 10 * time.Second}

type JobState int

type JobHandler func(ctx context.Context) error

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

func NewJob(handler JobHandler) *job {
	j := job{
		config: JobConfig{
			MaxTimeOut: defaultMaxTimeOut,
			Retries:    defaultRetryTime,
		},
		handler:    handler,
		state:      StateInit,
		retryIndex: -1,
		stopChan:   make(chan bool),
	}

	return &j
}

func (j *job) Excute(ctx context.Context) error {
	j.state = StateRunning

	err := j.handler(ctx)
	if err != nil {
		j.state = StateFailed
		return err
	}
	j.state = StateCompleted
	return nil
}

func (j *job) Retry(ctx context.Context) error {
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
