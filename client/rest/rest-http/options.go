package http

import (
	"net/http"
	"time"

	"github.com/taluos/Malt/client/rest/rest-http/internal/interceptors"
)

type clientOptions struct {
	address      string
	timeout      time.Duration
	retryCount   int
	userAgent    string
	headers      map[string]string
	interceptors []interceptors.Interceptor
	transport    http.RoundTripper
}

type ClientOption func(*clientOptions)

func WithAddress(baseURL string) ClientOption {
	return func(c *clientOptions) {
		c.address = baseURL
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *clientOptions) {
		c.timeout = timeout
	}
}

func WithRetryCount(count int) ClientOption {
	return func(c *clientOptions) {
		c.retryCount = count
	}
}

func WithUserAgent(userAgent string) ClientOption {
	return func(c *clientOptions) {
		c.userAgent = userAgent
	}
}

func WithHeader(key, value string) ClientOption {
	return func(c *clientOptions) {
		c.headers[key] = value
	}
}

func WithInterceptor(interceptor ...interceptors.Interceptor) ClientOption {
	return func(c *clientOptions) {
		c.interceptors = append(c.interceptors, interceptor...)
	}
}

func WithTransport(transport http.RoundTripper) ClientOption {
	return func(c *clientOptions) {
		c.transport = transport
	}
}
