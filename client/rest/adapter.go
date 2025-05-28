package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	restfasthttp "github.com/taluos/Malt/client/rest/rest-fasthttp"
	resthttp "github.com/taluos/Malt/client/rest/rest-http"
)

// httpClient 基于标准HTTP库的客户端适配器
type httpClient struct {
	client *resthttp.Client
}

// httpResponse 标准HTTP响应适配器
type httpResponse struct {
	resp *http.Response
	body []byte
}

// fastHTTPClient 基于FastHTTP库的客户端适配器
type fastHTTPClient struct {
	client *restfasthttp.Client
}

// fastHTTPResponse FastHTTP响应适配器
type fastHTTPResponse struct {
	resp *restfasthttp.Response
}

// 确保实现了接口
var _ Client = (*httpClient)(nil)
var _ Response = (*httpResponse)(nil)
var _ Client = (*fastHTTPClient)(nil)
var _ Response = (*fastHTTPResponse)(nil)

// httpClient 实现
func (h *httpClient) Get(ctx context.Context, path string, opts ...RequestOption) (Response, error) {
	resp, err := h.client.Get(ctx, path, convertRequestOptions(opts...)...)
	if err != nil {
		return nil, err
	}
	return wrapHTTPResponse(resp.Response)
}

func (h *httpClient) Post(ctx context.Context, path string, body interface{}, opts ...RequestOption) (Response, error) {
	resp, err := h.client.Post(ctx, path, body, convertRequestOptions(opts...)...)
	if err != nil {
		return nil, err
	}
	return wrapHTTPResponse(resp.Response)
}

func (h *httpClient) Put(ctx context.Context, path string, body interface{}, opts ...RequestOption) (Response, error) {
	resp, err := h.client.Put(ctx, path, body, convertRequestOptions(opts...)...)
	if err != nil {
		return nil, err
	}
	return wrapHTTPResponse(resp.Response)
}

func (h *httpClient) Delete(ctx context.Context, path string, opts ...RequestOption) (Response, error) {
	resp, err := h.client.Delete(ctx, path, convertRequestOptions(opts...)...)
	if err != nil {
		return nil, err
	}
	return wrapHTTPResponse(resp.Response)
}

func (h *httpClient) Patch(ctx context.Context, path string, body interface{}, opts ...RequestOption) (Response, error) {
	resp, err := h.client.Patch(ctx, path, body, convertRequestOptions(opts...)...)
	if err != nil {
		return nil, err
	}
	return wrapHTTPResponse(resp.Response)
}

func (h *httpClient) Close() error {
	return h.client.Close()
}

// httpResponse 实现
func (r *httpResponse) StatusCode() int {
	return r.resp.StatusCode
}

func (r *httpResponse) Body() []byte {
	return r.body
}

func (r *httpResponse) Header(key string) string {
	return r.resp.Header.Get(key)
}

func (r *httpResponse) JSON(v interface{}) error {
	return json.Unmarshal(r.body, v)
}

func (r *httpResponse) String() string {
	return string(r.body)
}

func (r *httpResponse) Reader() io.Reader {
	return bytes.NewReader(r.body)
}

// fastHTTPClient 实现
func (f *fastHTTPClient) Get(ctx context.Context, path string, opts ...RequestOption) (Response, error) {
	resp, err := f.client.Get(ctx, path, convertFastHTTPRequestOptions(opts...)...)
	if err != nil {
		return nil, err
	}
	return &fastHTTPResponse{resp: resp}, nil
}

func (f *fastHTTPClient) Post(ctx context.Context, path string, body interface{}, opts ...RequestOption) (Response, error) {
	resp, err := f.client.Post(ctx, path, body, convertFastHTTPRequestOptions(opts...)...)
	if err != nil {
		return nil, err
	}
	return &fastHTTPResponse{resp: resp}, nil
}

func (f *fastHTTPClient) Put(ctx context.Context, path string, body interface{}, opts ...RequestOption) (Response, error) {
	resp, err := f.client.Put(ctx, path, body, convertFastHTTPRequestOptions(opts...)...)
	if err != nil {
		return nil, err
	}
	return &fastHTTPResponse{resp: resp}, nil
}

func (f *fastHTTPClient) Delete(ctx context.Context, path string, opts ...RequestOption) (Response, error) {
	resp, err := f.client.Delete(ctx, path, convertFastHTTPRequestOptions(opts...)...)
	if err != nil {
		return nil, err
	}
	return &fastHTTPResponse{resp: resp}, nil
}

func (f *fastHTTPClient) Patch(ctx context.Context, path string, body interface{}, opts ...RequestOption) (Response, error) {
	resp, err := f.client.Patch(ctx, path, body, convertFastHTTPRequestOptions(opts...)...)
	if err != nil {
		return nil, err
	}
	return &fastHTTPResponse{resp: resp}, nil
}

func (f *fastHTTPClient) Close() error {
	return f.client.Close()
}

// fastHTTPResponse 实现
func (r *fastHTTPResponse) StatusCode() int {
	return r.resp.StatusCode()
}

func (r *fastHTTPResponse) Body() []byte {
	return r.resp.Body()
}

func (r *fastHTTPResponse) Header(key string) string {
	return r.resp.Header(key)
}

func (r *fastHTTPResponse) JSON(v interface{}) error {
	return r.resp.JSON(v)
}

func (r *fastHTTPResponse) String() string {
	return r.resp.String()
}

func (r *fastHTTPResponse) Reader() io.Reader {
	return r.resp.Reader()
}

// 工厂函数
func newHTTPClient(baseURL string, opts ...ClientOption) (Client, error) {
	clientOpts := convertToHTTPOptions(opts...)
	client := resthttp.NewClient(baseURL, clientOpts...)
	return &httpClient{client: client}, nil
}

func newFastHTTPClient(baseURL string, opts ...ClientOption) (Client, error) {
	clientOpts := convertToFastHTTPOptions(opts...)
	client := restfasthttp.NewClient(baseURL, clientOpts...)
	return &fastHTTPClient{client: client}, nil
}

// 选项转换函数
func convertToHTTPOptions(opts ...ClientOption) []resthttp.ClientOption {
	var httpOpts []resthttp.ClientOption
	clientOpts := &ClientOptions{}

	// 先应用所有选项到同一个结构体
	for _, opt := range opts {
		opt(clientOpts)
	}

	// 然后转换为HTTP选项
	if clientOpts.Timeout > 0 {
		httpOpts = append(httpOpts, resthttp.WithTimeout(clientOpts.Timeout))
	}
	if clientOpts.RetryCount > 0 {
		httpOpts = append(httpOpts, resthttp.WithRetryCount(clientOpts.RetryCount))
	}
	if clientOpts.UserAgent != "" {
		httpOpts = append(httpOpts, resthttp.WithUserAgent(clientOpts.UserAgent))
	}
	for k, v := range clientOpts.Headers {
		httpOpts = append(httpOpts, resthttp.WithHeader(k, v))
	}

	return httpOpts
}

func convertToFastHTTPOptions(opts ...ClientOption) []restfasthttp.ClientOption {
	var fasthttpOpts []restfasthttp.ClientOption

	for _, opt := range opts {
		clientOpts := &ClientOptions{}
		opt(clientOpts)

		if clientOpts.Timeout > 0 {
			fasthttpOpts = append(fasthttpOpts, restfasthttp.WithTimeout(clientOpts.Timeout))
		}
		if clientOpts.RetryCount > 0 {
			fasthttpOpts = append(fasthttpOpts, restfasthttp.WithRetryCount(clientOpts.RetryCount))
		}
		if clientOpts.UserAgent != "" {
			fasthttpOpts = append(fasthttpOpts, restfasthttp.WithUserAgent(clientOpts.UserAgent))
		}
		for k, v := range clientOpts.Headers {
			fasthttpOpts = append(fasthttpOpts, restfasthttp.WithHeader(k, v))
		}
	}

	return fasthttpOpts
}

func convertRequestOptions(opts ...RequestOption) []resthttp.RequestOption {
	reqOpts := &RequestOptions{}
	for _, opt := range opts {
		opt(reqOpts)
	}

	var httpReqOpts []resthttp.RequestOption
	for k, v := range reqOpts.Headers {
		httpReqOpts = append(httpReqOpts, resthttp.WithRequestHeader(k, v))
	}
	for k, v := range reqOpts.QueryParams {
		httpReqOpts = append(httpReqOpts, resthttp.WithQueryParam(k, v))
	}

	return httpReqOpts
}

func convertFastHTTPRequestOptions(opts ...RequestOption) []restfasthttp.RequestOption {
	reqOpts := &RequestOptions{}
	for _, opt := range opts {
		opt(reqOpts)
	}

	var fasthttpReqOpts []restfasthttp.RequestOption
	for k, v := range reqOpts.Headers {
		fasthttpReqOpts = append(fasthttpReqOpts, restfasthttp.WithRequestHeader(k, v))
	}
	for k, v := range reqOpts.QueryParams {
		fasthttpReqOpts = append(fasthttpReqOpts, restfasthttp.WithQueryParam(k, v))
	}

	return fasthttpReqOpts
}

func wrapHTTPResponse(resp *http.Response) (Response, error) {
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close() // 读取完成后再关闭
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	return &httpResponse{
		resp: resp,
		body: body,
	}, nil
}
