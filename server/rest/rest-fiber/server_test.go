package fiber

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	maltAgent "github.com/taluos/Malt/core/trace"
	"github.com/taluos/Malt/server/rest/rest-fiber/internal/auth"

	fiber "github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServerCreation 测试服务器创建
func TestServerCreation(t *testing.T) {
	tests := []struct {
		name string
		opts []ServerOptions
		want string
	}{
		{
			name: "default server",
			opts: nil,
			want: defaultName,
		},
		{
			name: "custom name server",
			opts: []ServerOptions{WithName("test-server")},
			want: "test-server",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewServer(tt.opts...)
			assert.NotNil(t, server)
			assert.NotNil(t, server.App)
			assert.Equal(t, tt.want, server.opts.name)
		})
	}
}

// TestServerOptions 测试服务器配置选项
func TestServerOptions(t *testing.T) {
	server := NewServer(
		WithName("test-server"),
		WithAddress(":9090"),
		WithTrans("en"),
		WithHealthz(false),
		WithEnableProfiling(false),
		WithEnableMetrics(true),
		WithEnableTracing(false),
		WithTrustedProxies([]string{"127.0.0.1"}),
	)

	assert.Equal(t, "test-server", server.opts.name)
	assert.Equal(t, ":9090", server.opts.address)
	assert.Equal(t, "en", server.opts.trans)
	assert.False(t, server.opts.enableHealth)
	assert.False(t, server.opts.enableProfiling)
	assert.True(t, server.opts.enableMetrics)
	assert.False(t, server.opts.enableTracing)
	assert.Equal(t, []string{"127.0.0.1"}, server.opts.trustedProxies)
}

// TestServerStartStop 测试服务器启动和关闭
func TestServerStartStop(t *testing.T) {
	server := NewServer(
		WithAddress(":0"), // 使用随机端口
		WithHealthz(true),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 启动服务器
	go func() {
		err := server.Start(ctx)
		if err != nil && !strings.Contains(err.Error(), "context canceled") {
			t.Errorf("Server start failed: %v", err)
		}
	}()

	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)

	// 测试健康检查
	resp, err := http.Get("http://127.0.0.1:8080/health")
	if err == nil {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, "ok", string(body))
	}

	// 关闭服务器
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer stopCancel()

	err = server.Stop(stopCtx)
	assert.NoError(t, err)
}

// TestMiddlewareApplication 测试中间件应用
func TestMiddlewareApplication(t *testing.T) {
	middlewareCalled := false
	testMiddleware := func(c fiber.Ctx) error {
		middlewareCalled = true
		return c.Next()
	}

	server := NewServer(
		WithMiddleware(testMiddleware),
		WithHealthz(true),
	)

	// 创建测试路由
	server.Get("/test", func(c fiber.Ctx) error {
		return c.SendString("test")
	})

	// 创建测试请求
	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/test"},
		Header: make(http.Header),
	}

	// 使用 Fiber 的测试功能
	resp, err := server.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, "test", string(body))
	assert.True(t, middlewareCalled)
}

// TestBasicAuthMiddleware 测试基础认证中间件
func TestBasicAuthMiddleware(t *testing.T) {
	// 创建基础认证策略
	basicStrategy := auth.NewBasicStrategy(func(username, password string) bool {
		return username == "admin" && password == "password"
	})

	authOperator := &auth.AuthOperator{}
	authOperator.SetStrategy(basicStrategy)

	server := NewServer(
		WithAuthOperator(authOperator),
	)

	// 创建受保护的路由
	server.Get("/protected", func(c fiber.Ctx) error {
		return c.SendString("protected resource")
	})

	tests := []struct {
		name           string
		authorization  string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "no authorization",
			authorization:  "",
			expectedStatus: 401,
		},
		{
			name:           "invalid authorization format",
			authorization:  "Invalid header",
			expectedStatus: 401,
		},
		{
			name:           "invalid credentials",
			authorization:  "Basic " + base64.StdEncoding.EncodeToString([]byte("wrong:credentials")),
			expectedStatus: 401,
		},
		{
			name:           "valid credentials",
			authorization:  "Basic " + base64.StdEncoding.EncodeToString([]byte("admin:password")),
			expectedStatus: 200,
			expectedBody:   "protected resource",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: "/protected"},
				Header: make(http.Header),
			}

			if tt.authorization != "" {
				req.Header.Set("Authorization", tt.authorization)
			}

			resp, err := server.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedBody != "" {
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedBody, string(body))
			}
		})
	}
}

// TestTracingMiddleware 测试追踪中间件
func TestTracingMiddleware(t *testing.T) {
	// 创建模拟的追踪代理
	agent := &maltAgent.Agent{}

	server := NewServer(
		WithEnableTracing(true),
		WithAgent(agent),
	)

	// 创建测试路由
	server.Get("/traced", func(c fiber.Ctx) error {
		return c.SendString("traced endpoint")
	})

	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/traced"},
		Header: make(http.Header),
	}

	// 添加追踪头
	req.Header.Set("traceparent", "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")

	resp, err := server.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "traced endpoint", string(body))
}

// TestHealthCheck 测试健康检查
func TestHealthCheck(t *testing.T) {
	tests := []struct {
		name         string
		enableHealth bool
		expectedCode int
		expectedBody string
	}{
		{
			name:         "health enabled",
			enableHealth: true,
			expectedCode: 200,
			expectedBody: "ok",
		},
		{
			name:         "health disabled",
			enableHealth: false,
			expectedCode: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewServer(
				WithHealthz(tt.enableHealth),
			)

			req := &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: "/health"},
				Header: make(http.Header),
			}

			resp, err := server.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			if tt.expectedBody != "" {
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedBody, string(body))
			}
		})
	}
}

// TestPProfMiddleware 测试性能分析中间件
func TestPProfMiddleware(t *testing.T) {
	server := NewServer(
		WithEnableProfiling(true),
	)

	// 测试pprof端点
	pprofEndpoints := []string{
		"/debug/pprof/",
		"/debug/pprof/heap",
		"/debug/pprof/goroutine",
		"/debug/pprof/allocs",
		"/debug/pprof/block",
		"/debug/pprof/cmdline",
		"/debug/pprof/mutex",
		"/debug/pprof/profile",
		"/debug/pprof/threadcreate",
		"/debug/pprof/trace",
	}

	for _, endpoint := range pprofEndpoints {
		t.Run(fmt.Sprintf("pprof_%s", endpoint), func(t *testing.T) {
			req := &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: endpoint},
				Header: make(http.Header),
			}

			resp, err := server.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// pprof端点应该返回200或者其他有效状态码
			assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 500)
		})
	}
}

// TestConcurrentRequests 测试并发请求
func TestConcurrentRequests(t *testing.T) {
	server := NewServer()

	// 创建测试路由
	server.Get("/concurrent", func(c fiber.Ctx) error {
		time.Sleep(10 * time.Millisecond) // 模拟处理时间
		return c.SendString("concurrent response")
	})

	const numRequests = 10
	results := make(chan error, numRequests)

	// 发送并发请求
	for i := 0; i < numRequests; i++ {
		go func() {
			req := &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: "/concurrent"},
				Header: make(http.Header),
			}

			resp, err := server.Test(req)
			if err != nil {
				results <- err
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				results <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				results <- err
				return
			}

			if string(body) != "concurrent response" {
				results <- fmt.Errorf("unexpected response body: %s", string(body))
				return
			}

			results <- nil
		}()
	}

	// 等待所有请求完成
	for i := 0; i < numRequests; i++ {
		err := <-results
		assert.NoError(t, err)
	}
}
