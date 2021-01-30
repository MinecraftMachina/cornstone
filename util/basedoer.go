package util

import (
	"github.com/ViRb3/sling/v2"
	"net/http"
	"net/http/cookiejar"
)

type BaseDoer struct {
	*http.Client
}

func (b *BaseDoer) Do(req *http.Request) (*http.Response, error) {
	response, err := b.Client.Do(req)
	if err != nil {
		return nil, err
	}
	for _, code := range ExpectedStatusCode {
		if response.StatusCode == code {
			return response, nil
		}
	}
	return nil, CreateUnexpectedStatusCodeError(req.RequestURI, response)
}

func newBaseClient() *http.Client {
	return &http.Client{}
}

func NewBaseDoer() (sling.Doer, *http.Client) {
	client := newBaseClient()
	doer := BaseDoer{Client: client}
	return &doer, client
}

func NewJarDoer() (sling.Doer, *http.Client) {
	jar, _ := cookiejar.New(nil)
	client := newBaseClient()
	client.Jar = jar
	doer := BaseDoer{Client: client}
	return &doer, client
}
