package rest

import (
	"bytes"
	"context"
	"io"

	restfasthttp "github.com/taluos/Malt/client/rest/rest-fasthttp"
)

// fastHTTPClient 基于FastHTTP库的客户端适配器
type fastHTTPClient struct {
	client *restfasthttp.Client
}

// fastHTTPResponse FastHTTP响应适配器
type fastHTTPResponse struct {
	resp *restfasthttp.Response
	body []byte
}

var _ Client = (*fastHTTPClient)(nil)
var _ Response = (*fastHTTPResponse)(nil)

func newFastHTTPClient(baseURL string, opts ...ClientOption) (Client, error) {
	clientOpts := convertToFastHTTPOptions(opts...)
	client := restfasthttp.NewClient(baseURL, clientOpts...)
	return &fastHTTPClient{client: client}, nil
}

// fastHTTPClient 实现
func (f *fastHTTPClient) Get(ctx context.Context, path string, opts ...RequestOption) (Response, error) {
	resp, err := f.client.Get(ctx, path, convertFastHTTPRequestOptions(opts...)...)
	if err != nil {
		return nil, err
	}
	return wrapFastHTTPResponse(resp)
}

func (f *fastHTTPClient) Post(ctx context.Context, path string, body interface{}, opts ...RequestOption) (Response, error) {
	resp, err := f.client.Post(ctx, path, body, convertFastHTTPRequestOptions(opts...)...)
	if err != nil {
		return nil, err
	}
	return wrapFastHTTPResponse(resp)
}

func (f *fastHTTPClient) Put(ctx context.Context, path string, body interface{}, opts ...RequestOption) (Response, error) {
	resp, err := f.client.Put(ctx, path, body, convertFastHTTPRequestOptions(opts...)...)
	if err != nil {
		return nil, err
	}
	return wrapFastHTTPResponse(resp)
}

func (f *fastHTTPClient) Delete(ctx context.Context, path string, opts ...RequestOption) (Response, error) {
	resp, err := f.client.Delete(ctx, path, convertFastHTTPRequestOptions(opts...)...)
	if err != nil {
		return nil, err
	}
	return wrapFastHTTPResponse(resp)
}

func (f *fastHTTPClient) Patch(ctx context.Context, path string, body interface{}, opts ...RequestOption) (Response, error) {
	resp, err := f.client.Patch(ctx, path, body, convertFastHTTPRequestOptions(opts...)...)
	if err != nil {
		return nil, err
	}
	return wrapFastHTTPResponse(resp)
}

func (f *fastHTTPClient) Close(_ context.Context) error {
	return f.client.Close()
}

// fastHTTPResponse 实现
func (r *fastHTTPResponse) StatusCode() int {
	return r.resp.StatusCode()
}

func (r *fastHTTPResponse) Body() []byte {
	return r.body
}

func (r *fastHTTPResponse) Header(key string) string {
	return r.resp.Header(key)
}

func (r *fastHTTPResponse) JSON(v interface{}) error {
	return r.resp.JSON(v)
}

func (r *fastHTTPResponse) String() string {
	return string(r.body)
}

func (r *fastHTTPResponse) Reader() io.Reader {
	return bytes.NewReader(r.body)
}

func convertToFastHTTPOptions(opts ...ClientOption) []restfasthttp.ClientOption {
	clientOpts := make([]restfasthttp.ClientOption, 0, len(opts))
	for _, opt := range opts {
		if co, ok := opt.(restfasthttp.ClientOption); ok {
			clientOpts = append(clientOpts, co)
		}
	}
	return clientOpts
}

func convertFastHTTPRequestOptions(opts ...RequestOption) []restfasthttp.RequestOption {
	requestOpts := make([]restfasthttp.RequestOption, 0, len(opts))
	for _, opt := range opts {
		if co, ok := opt.(restfasthttp.RequestOption); ok {
			requestOpts = append(requestOpts, co)
		}
	}
	return requestOpts
}

// 添加一个通用的响应包装函数
func wrapFastHTTPResponse(resp *restfasthttp.Response) (Response, error) {
	// 读取响应体
	body, err := io.ReadAll(resp.Reader())
	if err != nil {
		return nil, err
	}
	return &fastHTTPResponse{resp: resp, body: body}, nil
}
