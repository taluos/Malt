package interceptors

import (
	"context"
	"net/http"
)

type Interceptor interface {
	Intercept(ctx context.Context, req *http.Request, next RoundTripper) (*http.Response, error)
}

type RoundTripper func(ctx context.Context, req *http.Request) (*http.Response, error)

// 认证拦截器
type AuthInterceptor struct {
	token string
}

func NewAuthInterceptor(token string) *AuthInterceptor {
	return &AuthInterceptor{token: token}
}

func (a *AuthInterceptor) Intercept(ctx context.Context, req *http.Request, next RoundTripper) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+a.token)
	return next(ctx, req)
}
