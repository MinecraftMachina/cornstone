package util

import (
	"context"
	"github.com/cavaliercoder/grab"
	"time"
)

const SizeDivisor = 1_000_000 // 1MB

type MultiDownloader struct {
	requests []*grab.Request
	workers  int
	client   *grab.Client
}

func NewMultiDownloader(workers int, requests ...*grab.Request) *MultiDownloader {
	downloader := MultiDownloader{
		requests: requests,
		workers:  workers,
		client:   grab.NewClient(),
	}
	downloader.client.UserAgent = DefaultUserAgent
	return &downloader
}

func (s *MultiDownloader) Do() (<-chan *grab.Response, context.CancelFunc) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	for i := range s.requests {
		s.requests[i] = s.requests[i].WithContext(ctx)
	}

	clientResult := s.client.DoBatch(s.workers, s.requests...)
	returnResult := make(chan *grab.Response, s.workers)

	if len(s.requests) > 1 {
		go s.handleMultiFile(clientResult, returnResult)
	} else {
		go s.handleSingleFile(clientResult, returnResult)
	}
	return returnResult, cancelFunc
}

func (s *MultiDownloader) handleMultiFile(result <-chan *grab.Response, returnResult chan<- *grab.Response) {
	defer close(returnResult)
	bar := NewBar(len(s.requests))
	defer bar.Finish()
	for response := range result {
		returnResult <- response
		bar.Add(1)
	}
}

func (s *MultiDownloader) handleSingleFile(result <-chan *grab.Response, returnResult chan<- *grab.Response) {
	defer close(returnResult)
	response := <-result
	barSize := int(response.Size() / SizeDivisor)
	if barSize < 1 {
		// spinner mode
		barSize = -1
	}

	t := time.NewTicker(200 * time.Millisecond)
	defer t.Stop()

	bar := NewBar(barSize)
	defer bar.Finish()

	var lastSize int
Loop:
	for {
		select {
		case <-t.C:
			delta := int(response.BytesComplete()/SizeDivisor) - lastSize
			// set will break spinner above 10, ref: progressbar.go#L400
			bar.Add(delta)
			lastSize += delta
		case <-response.Done:
			break Loop
		}
	}
	returnResult <- response
}
