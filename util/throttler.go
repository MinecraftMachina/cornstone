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
				wg.Done()
			} else {
				throttleChan <- true
				go func() {
					sourceItem := t.Source[i]
					result, err := t.Operation(sourceItem)
					resultChan <- Result{Data: result, Error: err}
					bar.Add(1)
					wg.Done()
					<-throttleChan
				}()
			}
		}
	}()

	go func() {
		wg.Wait()
		bar.Finish()
		close(resultChan)
	}()

	return resultChan
}
