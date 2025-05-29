package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

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

// 确保实现了接口
var _ Client = (*httpClient)(nil)
var _ Response = (*httpResponse)(nil)

// 工厂函数
func newHTTPClient(baseURL string, opts ...ClientOption) (Client, error) {
	clientOpts := convertToHTTPOptions(opts...)
	client := resthttp.NewClient(baseURL, clientOpts...)
	return &httpClient{client: client}, nil
}

func (h *httpClient) Get(ctx context.Context, path string, opts ...RequestOption) (Response, error) {
	resp, err := h.client.Get(ctx, path, convertRequestOptions(opts...)...)
	if err != nil {
		return nil, err
	}
	return wrapHTTPResponse(resp)
}

func (h *httpClient) Post(ctx context.Context, path string, body interface{}, opts ...RequestOption) (Response, error) {
	resp, err := h.client.Post(ctx, path, body, convertRequestOptions(opts...)...)
	if err != nil {
		return nil, err
	}
	return wrapHTTPResponse(resp)
}

func (h *httpClient) Put(ctx context.Context, path string, body interface{}, opts ...RequestOption) (Response, error) {
	resp, err := h.client.Put(ctx, path, body, convertRequestOptions(opts...)...)
	if err != nil {
		return nil, err
	}
	return wrapHTTPResponse(resp)
}

func (h *httpClient) Delete(ctx context.Context, path string, opts ...RequestOption) (Response, error) {
	resp, err := h.client.Delete(ctx, path, convertRequestOptions(opts...)...)
	if err != nil {
		return nil, err
	}
	return wrapHTTPResponse(resp)
}

func (h *httpClient) Patch(ctx context.Context, path string, body interface{}, opts ...RequestOption) (Response, error) {
	resp, err := h.client.Patch(ctx, path, body, convertRequestOptions(opts...)...)
	if err != nil {
		return nil, err
	}
	return wrapHTTPResponse(resp)
}

func (h *httpClient) Close(ctx context.Context) error {
	return h.client.Close(ctx)
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

func convertToHTTPOptions(opts ...ClientOption) []resthttp.ClientOption {
	clientOption := make([]resthttp.ClientOption, 0, len(opts))
	for _, opt := range opts {
		if co, ok := opt.(resthttp.ClientOption); ok {
			clientOption = append(clientOption, co)
		}
	}
	return clientOption
}

func convertRequestOptions(opts ...RequestOption) []resthttp.RequestOption {
	requestOption := make([]resthttp.RequestOption, 0, len(opts))
	for _, opt := range opts {
		if ro, ok := opt.(resthttp.RequestOption); ok {
			requestOption = append(requestOption, ro)
		}
	}
	return requestOption
}

// 添加一个通用的响应包装函数
func wrapHTTPResponse(resp *resthttp.Response) (Response, error) {
	// 读取响应体
	body, err := io.ReadAll(resp.Response.Body)
	if err != nil {
		return nil, err
	}
	resp.Response.Body.Close() // 关闭原始body
	return &httpResponse{resp: resp.Response, body: body}, nil
}
