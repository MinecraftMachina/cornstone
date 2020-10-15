package util

import (
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
	result := s.client.DoBatch(s.workers, s.requests...)

	if len(s.requests) < 2 {
		return s.handleSingleFile(result)
	} else {
		return s.handleMultiFile(result)
	}
}

func (s *MultiDownloader) handleMultiFile(result <-chan *grab.Response) error {
	bar := NewBar(len(s.requests))
	defer bar.Finish()
	for response := range result {
		if err := response.Err(); err != nil {
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
	if err := response.Err(); err != nil {
		return err
	}
	return nil
}
