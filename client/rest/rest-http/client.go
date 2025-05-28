package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/taluos/Malt/client/rest/rest-http/internal/interceptors"
)

type Client struct {
	*http.Client
	opts *clientOptions
}

func NewClient(baseURL string, opts ...ClientOption) *Client {
	o := &clientOptions{
		address:    baseURL, // 设置 baseURL
		timeout:    defaultTimeout,
		retryCount: defaultCount,
		userAgent:  defaultAgent,
		headers:    make(map[string]string),
	}

	for _, opt := range opts {
		opt(o)
	}

	transport := o.transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	cli := &Client{&http.Client{Timeout: o.timeout, Transport: transport}, o}

	return cli
}

func (c *Client) Get(ctx context.Context, path string, opts ...RequestOption) (*Response, error) {
	return c.doRequest(ctx, http.MethodGet, path, nil, opts...)
}

func (c *Client) Post(ctx context.Context, path string, body interface{}, opts ...RequestOption) (*Response, error) {
	return c.doRequest(ctx, http.MethodPost, path, body, opts...)
}

func (c *Client) Put(ctx context.Context, path string, body interface{}, opts ...RequestOption) (*Response, error) {
	return c.doRequest(ctx, http.MethodPut, path, body, opts...)
}

func (c *Client) Delete(ctx context.Context, path string, opts ...RequestOption) (*Response, error) {
	return c.doRequest(ctx, http.MethodDelete, path, nil, opts...)
}

func (c *Client) Patch(ctx context.Context, path string, body interface{}, opts ...RequestOption) (*Response, error) {
	return c.doRequest(ctx, http.MethodPatch, path, body, opts...)
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, opts ...RequestOption) (*Response, error) {
	reqOpts := &requestOptions{}
	for _, opt := range opts {
		opt(reqOpts)
	}

	// 构建URL
	fullURL := c.opts.address + "/" + strings.TrimPrefix(path, "/")
	if len(reqOpts.queryParams) > 0 {
		u, err := url.Parse(fullURL)
		if err != nil {
			return nil, fmt.Errorf("invalid URL: %w", err)
		}
		q := u.Query()
		for k, v := range reqOpts.queryParams {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		fullURL = u.String()
	}

	// 处理请求体
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// 设置默认头部
	req.Header.Set("User-Agent", c.opts.userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// 设置全局头部
	for k, v := range c.opts.headers {
		req.Header.Set(k, v)
	}

	// 设置请求级头部
	for k, v := range reqOpts.headers {
		req.Header.Set(k, v)
	}

	// 执行拦截器链
	return c.executeWithInterceptors(ctx, req)
}

func (c *Client) executeWithInterceptors(ctx context.Context, req *http.Request) (*Response, error) {
	if len(c.opts.interceptors) == 0 {
		res, err := c.Do(req)
		response := NewResponse(res)
		return response, err
	}

	// 构建拦截器链
	var handler interceptors.RoundTripper
	handler = func(ctx context.Context, req *http.Request) (*http.Response, error) {
		return c.Do(req)
	}

	// 反向遍历拦截器，构建调用链
	for i := len(c.opts.interceptors) - 1; i >= 0; i-- {
		interceptor := c.opts.interceptors[i]
		next := handler
		handler = func(ctx context.Context, req *http.Request) (*http.Response, error) {
			return interceptor.Intercept(ctx, req, next)
		}
	}
	res, err := handler(ctx, req)
	response := NewResponse(res)
	return response, err
}

func (c *Client) Close() error {
	return c.Close()
}
