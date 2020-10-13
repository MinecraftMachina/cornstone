package util

import (
	"fmt"
	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

const DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36"

var defaultHeaders = map[string]string{
	"User-Agent": DefaultUserAgent,
}

var (
	ErrNon200StatusCode = errors.New("non-200 status code")
)

func CreateNon200Error(code int, body []byte) error {
	return errors.WithMessage(errors.WithStack(ErrNon200StatusCode),
		fmt.Sprintf("code: %d, body: %s", code, string(body)))
}

var DefaultClient *sling.Sling

func init() {
	doer, _ := NewBaseDoer()
	DefaultClient = sling.New().Doer(doer).SetMany(defaultHeaders)
}
