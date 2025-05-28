package fasthttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

type Client struct {
	client  *fasthttp.Client
	baseURL string
	opts    *clientOptions
}

func NewClient(baseURL string, opts ...ClientOption) *Client {
	o := &clientOptions{
		timeout:             30 * time.Second,
		retryCount:          3,
		userAgent:           "Malt-FastHTTP-Client",
		headers:             make(map[string]string),
		maxConnsPerHost:     512,
		maxIdleConnDuration: 10 * time.Second,
		readTimeout:         10 * time.Second,
		writeTimeout:        10 * time.Second,
	}

	for _, opt := range opts {
		opt(o)
	}

	client := &fasthttp.Client{
		MaxConnsPerHost:     o.maxConnsPerHost,
		MaxIdleConnDuration: o.maxIdleConnDuration,
		ReadTimeout:         o.readTimeout,
		WriteTimeout:        o.writeTimeout,
	}

	return &Client{
		client:  client,
		baseURL: strings.TrimSuffix(baseURL, "/"),
		opts:    o,
	}
}

func (c *Client) Get(ctx context.Context, path string, opts ...RequestOption) (*Response, error) {
	return c.doRequest(ctx, fasthttp.MethodGet, path, nil, opts...)
}

func (c *Client) Post(ctx context.Context, path string, body interface{}, opts ...RequestOption) (*Response, error) {
	return c.doRequest(ctx, fasthttp.MethodPost, path, body, opts...)
}

func (c *Client) Put(ctx context.Context, path string, body interface{}, opts ...RequestOption) (*Response, error) {
	return c.doRequest(ctx, fasthttp.MethodPut, path, body, opts...)
}

func (c *Client) Delete(ctx context.Context, path string, opts ...RequestOption) (*Response, error) {
	return c.doRequest(ctx, fasthttp.MethodDelete, path, nil, opts...)
}

func (c *Client) Patch(ctx context.Context, path string, body interface{}, opts ...RequestOption) (*Response, error) {
	return c.doRequest(ctx, fasthttp.MethodPatch, path, body, opts...)
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, opts ...RequestOption) (*Response, error) {
	reqOpts := &requestOptions{}
	for _, opt := range opts {
		opt(reqOpts)
	}

	// 构建URL
	fullURL := c.baseURL + "/" + strings.TrimPrefix(path, "/")
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

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	// 设置请求方法和URL
	req.Header.SetMethod(method)
	req.SetRequestURI(fullURL)

	// 设置默认头部
	req.Header.SetUserAgent(c.opts.userAgent)

	// 设置全局头部
	for k, v := range c.opts.headers {
		req.Header.Set(k, v)
	}

	// 设置请求级头部
	for k, v := range reqOpts.headers {
		req.Header.Set(k, v)
	}

	// 处理请求体
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		req.SetBody(bodyBytes)
		req.Header.SetContentType("application/json")
	}

	// 执行请求
	err := c.client.DoTimeout(req, resp, c.opts.timeout)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	// 复制响应数据
	response := &Response{
		Response: &fasthttp.Response{},
	}
	resp.CopyTo(response.Response)

	return response, nil
}

func (c *Client) Close() error {
	// FastHTTP客户端通常不需要显式关闭
	return nil
}
