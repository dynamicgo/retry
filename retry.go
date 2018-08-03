package retry

import (
	"context"
	"sync"
	"time"
)

// Action retry action object
type Action interface {
	Do() <-chan bool // do action with retry strategy
	Error() <-chan error
	Close() // close action and cancel retry option
}

// BackoffStrategy .
type BackoffStrategy interface {
	Next(ctx context.Context) bool
}

// TimesStrategy .
type TimesStrategy interface {
	Next(times int) bool
}

type actionConfig struct {
}

type actionImpl struct {
	f       func() error
	backoff BackoffStrategy
	times   TimesStrategy
	done    chan bool
	errors  chan error
	ctx     context.Context
	cancel  context.CancelFunc
	once    sync.Once
}

// Option .
type Option func(action *actionImpl)

// New create a new action object for retry
func New(f func() error, options ...Option) Action {
	return initWithOptions(context.Background(), f, options...)
}

// NewWitContext create a new action object for retry with ctx
func NewWitContext(ctx context.Context, f func() error, options ...Option) Action {
	return initWithOptions(ctx, f, options...)
}

func initWithOptions(ctx context.Context, f func() error, options ...Option) *actionImpl {

	ctx2, cancel := context.WithCancel(context.Background())

	action := &actionImpl{
		f:      f,
		done:   make(chan bool),
		errors: make(chan error),
		ctx:    ctx2,
		cancel: cancel,
		backoff: &backoffStrategy{
			duration: time.Second,
			backoff:  2.0,
		},
		times: &timesStrategy{
			maxtimes: 10,
		},
	}

	for _, option := range options {
		option(action)
	}

	return action
}

func (action *actionImpl) Do() <-chan bool {

	go action.once.Do(action.doAction)

	return action.done
}

func (action *actionImpl) doAction() {
	for i := 0; action.times.Next(i); i++ {
		if err := action.f(); err != nil {
			action.errors <- err
			if action.backoff.Next(action.ctx) {
				continue
			}
		}

		action.done <- true
	}

	action.done <- false
}

func (action *actionImpl) Close() {
	action.cancel()
}

func (action *actionImpl) Error() <-chan error {
	return action.errors
}
