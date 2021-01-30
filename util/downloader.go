package util

import (
	"context"
	"github.com/ViRb3/go-parallel/downloader"
	"github.com/pkg/errors"
)

func MultiDownload(config downloader.SharedConfig) (<-chan downloader.Result, context.CancelFunc) {
	returnResult := make(chan downloader.Result, 10)
	config.ExpectedStatusCodes = append(config.ExpectedStatusCodes, 200)
	results, cancelFunc := downloader.NewMultiDownloaderSling(downloader.ConfigSling{
		Client:       defaultClientNoStatusCheck,
		SharedConfig: config,
	}).Run()
	go func() {
		for result := range results {
			if errors.Is(result.Err, downloader.ErrUnexpectedStatusCode) {
				result.Err = CreateUnexpectedStatusCodeError(result.Job.Url, result.Response)
			}
			returnResult <- result
		}
		close(returnResult)
	}()
	return returnResult, cancelFunc
}
