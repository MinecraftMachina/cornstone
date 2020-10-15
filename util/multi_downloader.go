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

func (s *MultiDownloader) Do() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	for i := range s.requests {
		s.requests[i] = s.requests[i].WithContext(ctx)
	}

	result := s.client.DoBatch(s.workers, s.requests...)

	if len(s.requests) > 1 {
		return s.handleMultiFile(result, cancelFunc)
	} else {
		return s.handleSingleFile(result)
	}
}

func (s *MultiDownloader) handleMultiFile(result <-chan *grab.Response, cancelFunc context.CancelFunc) error {
	bar := NewBar(len(s.requests))
	defer bar.Finish()
	for response := range result {
		if err := response.Err(); err != nil {
			cancelFunc()
			return err
		}
		bar.Add(1)
	}
	return nil
}

func (s *MultiDownloader) handleSingleFile(result <-chan *grab.Response) error {
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

Loop:
	for {
		select {
		case <-t.C:
			bar.Set(int(response.BytesComplete() / SizeDivisor))
		case <-response.Done:
			break Loop
		}
	}
	return response.Err()
}
