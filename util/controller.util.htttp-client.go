package util

import (
	"net/http"
	"time"
)

// bare bone client, can be good if you want full control over each request
func SharedHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 45 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        200,
			MaxIdleConnsPerHost: 50,
			MaxConnsPerHost:     100,
			IdleConnTimeout:     90 * time.Second,
			DisableKeepAlives:   false,
			DisableCompression:  false,
		},
	}
}

// Q. What is the use case of this file?
// A. Using this file to make like http client to make get post requests,
// since we make a lot of requests to huggingface api, so to optimize it we use this.

// --------------------------------------------------------------------------------------------//

// we can also use resty, automatic retries, timeouts, connection pooling etc.
/*
package util

import (
	"github.com/go-resty/resty/v2"
	"time"
	)

	type RestyClient struct {
		client *resty.Client
		}

func SharedHTTPClient() *RestyClient {
	c := resty.New().
		SetTimeout(45 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetHeader("Accept-Encoding", "gzip")

	return &RestyClient{client: c}
}

func (r *RestyClient) Get(url string) (*resty.Response, error) {
	return r.client.R().Get(url)
}

func (r *RestyClient) Post(url string, body interface{}) (*resty.Response, error) {
	return r.client.R().SetBody(body).Post(url)
}
*/
