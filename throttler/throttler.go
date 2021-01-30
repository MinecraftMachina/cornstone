package throttler

import (
	"context"
	"cornstone/util"
	"sync"
)

type Throttler struct {
	Config
}

type Config struct {
	Ctx          context.Context
	ResultBuffer int
	Workers      int
	Source       []interface{}
	Operation    func(sourceItem interface{}) interface{}
}

func NewThrottler(config Config) *Throttler {
	return &Throttler{
		config,
	}
}

func (t *Throttler) Run() <-chan interface{} {
	throttleChan := make(chan bool, t.Workers)
	resultChan := make(chan interface{}, t.ResultBuffer)
	wg := sync.WaitGroup{}

	bar := util.NewBar(len(t.Source))
	wg.Add(len(t.Source))

	go func() {
		for i := range t.Source {
			if t.Ctx.Err() != nil {
				wg.Done()
			} else {
				sourceItem := t.Source[i]
				throttleChan <- true
				go func() {
					result := t.Operation(sourceItem)
					resultChan <- result
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
