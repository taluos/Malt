package client

import (
	"context"

	"github.com/taluos/Malt/client/rest"
)

// restClientWrapper REST客户端包装器
type restClientWrapper struct {
	client     rest.Client
	clientType string
}

// RESTClient REST客户端接口，继承统一接口
type RESTClient interface {
	Type() string
	rest.Client // 嵌入REST客户端接口
}

var _ RESTClient = (*restClientWrapper)(nil)

// NewRESTClient 创建REST客户端
func newRESTClient(clientType string, baseURL string, opts ...rest.ClientOption) (RESTClient, error) {
	client, err := rest.NewClient(clientType, baseURL, opts...)
	if err != nil {
		return nil, err
	}
	return &restClientWrapper{client: client, clientType: clientType}, nil
}

func (w *restClientWrapper) Type() string {
	return w.clientType
}

func (w *restClientWrapper) Close(ctx context.Context) error {
	return w.client.Close(ctx)
}

// 嵌入REST客户端的所有方法
func (w *restClientWrapper) Get(ctx context.Context, path string, opts ...rest.RequestOption) (rest.Response, error) {
	return w.client.Get(ctx, path, opts...)
}

func (w *restClientWrapper) Post(ctx context.Context, path string, body any, opts ...rest.RequestOption) (rest.Response, error) {
	return w.client.Post(ctx, path, body, opts...)
}

func (w *restClientWrapper) Put(ctx context.Context, path string, body any, opts ...rest.RequestOption) (rest.Response, error) {
	return w.client.Put(ctx, path, body, opts...)
}

func (w *restClientWrapper) Delete(ctx context.Context, path string, opts ...rest.RequestOption) (rest.Response, error) {
	return w.client.Delete(ctx, path, opts...)
}

func (w *restClientWrapper) Patch(ctx context.Context, path string, body any, opts ...rest.RequestOption) (rest.Response, error) {
	return w.client.Patch(ctx, path, body, opts...)
}
