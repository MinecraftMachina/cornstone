package util

import (
	"context"
	"sync"
)

type Result struct {
	Data  interface{}
	Error error
}

type Throttler struct {
	ThrottlerConfig
}

type ThrottlerConfig struct {
	Ctx          context.Context
	ResultBuffer int
	Workers      int
	Source       []interface{}
	Operation    func(sourceItem interface{}) (interface{}, error)
}

func NewThrottler(config ThrottlerConfig) *Throttler {
	return &Throttler{
		config,
	}
}

func (t *Throttler) Run() <-chan Result {
	throttleChan := make(chan bool, t.Workers)
	resultChan := make(chan Result, t.ResultBuffer)
	wg := sync.WaitGroup{}

	bar := NewBar(len(t.Source))
	wg.Add(len(t.Source))

	go func() {
		for i := range t.Source {
			if t.Ctx.Err() != nil {
				return
			}
			sourceItem := t.Source[i]
			throttleChan <- true
			go func() {
				defer func() {
					wg.Done()
					<-throttleChan
					bar.Add(1)
				}()
				result, err := t.Operation(sourceItem)
				resultChan <- Result{Data: result, Error: err}
			}()
		}
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return resultChan
}
