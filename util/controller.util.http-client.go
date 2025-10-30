package util

import (
	"net/http"
	"time"
)

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
