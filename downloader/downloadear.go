package downloader

import (
	"context"
	"cornstone/throttler"
	"cornstone/util"
	"github.com/ViRb3/sling/v2"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"strconv"
)

type Request struct {
	DownloadPath string
	DownloadUrl  string
	Tag          interface{}
}

type Result struct {
	Response *http.Response
	Request  Request
	Err      error
}

type MultiDownloader struct {
	requests     []Request
	workers      int
	client       *sling.Sling
	skipExisting bool
}

func NewMultiDownloader(workers int, requests ...Request) *MultiDownloader {
	downloader := MultiDownloader{
		requests: requests,
		workers:  workers,
		client:   util.DefaultClient.New(),
	}
	return &downloader
}

func (s *MultiDownloader) Do() (<-chan Result, context.CancelFunc) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	var source []interface{}
	for i := range s.requests {
		source = append(source, s.requests[i])
	}

	downloadThrottler := throttler.NewThrottler(throttler.Config{
		Ctx:          ctx,
		ResultBuffer: 0,
		Workers:      s.workers,
		Source:       source,
		Operation: func(sourceItem interface{}) interface{} {
			request := sourceItem.(Request)

			resp, err := s.client.New().Head(request.DownloadUrl).Receive(nil, nil)
			// if HEAD is unsupported, don't panic
			if err != nil && !errors.Is(err, util.ErrNon200StatusCode) {
				return Result{resp, request, err}
			}
			// if HEAD is supported, compare file lengths in hopes to skip re-downloading
			if err == nil {
				if contentLen := resp.Header.Get("Content-Length"); contentLen != "" {
					contentLenInt, err := strconv.ParseInt(contentLen, 10, 64)
					if err == nil {
						if stat, err := os.Stat(request.DownloadPath); err == nil && stat.Size() == contentLenInt {
							return Result{resp, request, nil}
						}
					}
				}
			}
			resp, err = s.client.New().Get(request.DownloadUrl).ReceiveBody()
			if err != nil {
				return Result{resp, request, err}
			}
			defer resp.Body.Close()
			file, err := os.Create(request.DownloadPath)
			if err != nil {
				return Result{resp, request, err}
			}
			defer file.Close()
			if _, err := io.Copy(file, resp.Body); err != nil {
				os.Remove(file.Name())
				return Result{resp, request, err}
			}
			return Result{resp, request, nil}
		},
	})

	returnResult := make(chan Result, 10)
	go s.handleMultiFile(returnResult, downloadThrottler)
	return returnResult, cancelFunc
}

func (s *MultiDownloader) handleMultiFile(returnResult chan Result, downloadThrottler *throttler.Throttler) {
	defer close(returnResult)
	bar := util.NewBar(len(s.requests))
	defer bar.Finish()
	for result := range downloadThrottler.Run() {
		returnResult <- result.(Result)
		bar.Add(1)
	}
}
