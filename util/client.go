package util

import (
	"cornstone/aliases/e"
	"fmt"
	"github.com/ViRb3/go-parallel/downloader"
	"github.com/ViRb3/sling/v2"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

const DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36"

var (
	ErrUnexpectedStatusCode = downloader.ErrUnexpectedStatusCode
	defaultHeaders          = map[string]string{
		"User-Agent": DefaultUserAgent,
	}
	ExpectedStatusCode = []int{200}
)

func CreateUnexpectedStatusCodeError(url string, resp *http.Response) error {
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return errors.WithMessagef(e.S(ErrUnexpectedStatusCode),
		fmt.Sprintf("url: %s, code: %d, body: %s", url, resp.StatusCode, string(bodyBytes)))
}

var DefaultClient *sling.Sling
var defaultClientNoStatusCheck *sling.Sling

func init() {
	doer, _ := NewBaseDoer()
	DefaultClient = sling.New().Doer(doer).SetMany(defaultHeaders)
	defaultClientNoStatusCheck = sling.New().SetMany(defaultHeaders)
}
