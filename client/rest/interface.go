package rest

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/taluos/Malt/pkg/errors"
)

// Client 定义了HTTP客户端的通用接口
type Client interface {
	// Get 执行GET请求
	Get(ctx context.Context, path string, opts ...RequestOption) (Response, error)
	// Post 执行POST请求
	Post(ctx context.Context, path string, body interface{}, opts ...RequestOption) (Response, error)
	// Put 执行PUT请求
	Put(ctx context.Context, path string, body interface{}, opts ...RequestOption) (Response, error)
	// Delete 执行DELETE请求
	Delete(ctx context.Context, path string, opts ...RequestOption) (Response, error)
	// Patch 执行PATCH请求
	Patch(ctx context.Context, path string, body interface{}, opts ...RequestOption) (Response, error)
	// Close 关闭客户端
	Close() error
}

// Response 定义了HTTP响应的通用接口
type Response interface {
	// StatusCode 返回HTTP状态码
	StatusCode() int
	// Body 返回响应体
	Body() []byte
	// Header 返回响应头
	Header(key string) string
	// JSON 将响应体解析为JSON
	JSON(v interface{}) error
	// String 返回响应体字符串
	String() string
	// Reader 返回响应体读取器
	Reader() io.Reader
}

// RequestOption 定义了请求选项的函数类型
type RequestOption func(*RequestOptions)

// RequestOptions 定义了请求选项结构
type RequestOptions struct {
	Headers     map[string]string
	QueryParams map[string]string
	Timeout     *time.Duration
}

// ClientOption 定义了客户端选项的函数类型
type ClientOption func(*ClientOptions)

// ClientOptions 定义了客户端选项结构
type ClientOptions struct {
	Timeout    time.Duration
	RetryCount int
	UserAgent  string
	Headers    map[string]string
	BaseURL    string
}

const (
	HTTPClient     string = "http"
	FastHTTPClient string = "fasthttp"
)

// NewClient 创建新的HTTP客户端
func NewClient(clientType string, baseURL string, opts ...ClientOption) (Client, error) {
	switch clientType {
	case HTTPClient:
		return newHTTPClient(baseURL, opts...)
	case FastHTTPClient:
		return newFastHTTPClient(baseURL, opts...)
	default:
		return nil, errors.New(fmt.Sprintf("unsupported client type: %s", clientType))
	}
}
