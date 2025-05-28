package fasthttp

import "time"

type clientOptions struct {
	timeout             time.Duration
	retryCount          int
	userAgent           string
	headers             map[string]string
	maxConnsPerHost     int
	maxIdleConnDuration time.Duration
	readTimeout         time.Duration
	writeTimeout        time.Duration
}

type ClientOption func(*clientOptions)

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

func WithMaxConnsPerHost(maxConns int) ClientOption {
	return func(c *clientOptions) {
		c.maxConnsPerHost = maxConns
	}
}

func WithMaxIdleConnDuration(duration time.Duration) ClientOption {
	return func(c *clientOptions) {
		c.maxIdleConnDuration = duration
	}
}

func WithReadTimeout(timeout time.Duration) ClientOption {
	return func(c *clientOptions) {
		c.readTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) ClientOption {
	return func(c *clientOptions) {
		c.writeTimeout = timeout
	}
}
