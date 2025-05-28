package httpserver

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name string
		opts []ServerOptions
		want func(*Server) bool
	}{
		{
			name: "default server",
			opts: nil,
			want: func(s *Server) bool {
				return s.opts.name == defaultName &&
					s.opts.address == defaultAddr &&
					s.opts.mode == gin.DebugMode &&
					s.opts.enableHealth == true &&
					s.opts.enableProfiling == true &&
					s.opts.enableMetrics == false
			},
		},
		{
			name: "custom server with options",
			opts: []ServerOptions{
				WithName("test-server"),
				WithAddress("localhost:9090"),
				WithMode(gin.ReleaseMode),
				WithEnableMetrics(true),
				WithHealthz(false),
			},
			want: func(s *Server) bool {
				return s.opts.name == "test-server" &&
					s.opts.address == "localhost:9090" &&
					s.opts.mode == gin.ReleaseMode &&
					s.opts.enableHealth == false &&
					s.opts.enableMetrics == true
			},
		},
		{
			name: "server with middleware",
			opts: []ServerOptions{
				WithMiddleware(gin.Logger(), gin.Recovery()),
			},
			want: func(s *Server) bool {
				return len(s.opts.middlewares) == 2
			},
		},
		{
			name: "server with trusted proxies",
			opts: []ServerOptions{
				WithTrustedProxies([]string{"127.0.0.1", "192.168.1.1"}),
			},
			want: func(s *Server) bool {
				return len(s.opts.trustedProxies) == 2 &&
					s.opts.trustedProxies[0] == "127.0.0.1" &&
					s.opts.trustedProxies[1] == "192.168.1.1"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewServer(tt.opts...)
			require.NotNil(t, server)
			require.NotNil(t, server.Engine)
			assert.True(t, tt.want(server))
		})
	}
}

func TestServerMethods(t *testing.T) {
	server := NewServer(
		WithName("test-server"),
		WithAddress("localhost:8081"),
	)

	// 测试 Name 方法
	assert.Equal(t, "test-server", server.Name())

	// 测试 Address 方法
	assert.Equal(t, "localhost:8081", server.Address())

	// 测试 Mode 方法
	assert.Equal(t, gin.DebugMode, server.Mode())

	// 测试 Trans 方法
	trans := server.Trans()
	assert.NotNil(t, trans)
}

func TestServerHealthEndpoint(t *testing.T) {
	server := NewServer(
		WithAddress("localhost:8082"),
		WithHealthz(true),
	)

	// 启动服务器
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := server.Start(ctx)
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("Server start failed: %v", err)
		}
	}()

	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)

	// 测试健康检查端点
	resp, err := http.Get("http://localhost:8082/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 停止服务器
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopCancel()
	err = server.Stop(stopCtx)
	assert.NoError(t, err)
}

func TestServerWithoutHealth(t *testing.T) {
	server := NewServer(
		WithAddress("localhost:8083"),
		WithHealthz(false),
	)

	// 启动服务器
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := server.Start(ctx)
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("Server start failed: %v", err)
		}
	}()

	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)

	// 测试健康检查端点应该不存在
	resp, err := http.Get("http://localhost:8083/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// 停止服务器
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopCancel()
	err = server.Stop(stopCtx)
	assert.NoError(t, err)
}

func TestServerStartStop(t *testing.T) {
	server := NewServer(
		WithAddress("localhost:8084"),
	)

	ctx := context.Background()

	// 测试启动服务器
	go func() {
		err := server.Start(ctx)
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("Server start failed: %v", err)
		}
	}()

	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)

	// 验证服务器正在运行
	resp, err := http.Get("http://localhost:8084/health")
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 测试停止服务器
	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Stop(stopCtx)
	assert.NoError(t, err)

	// 验证服务器已停止
	time.Sleep(100 * time.Millisecond)
	_, err = http.Get("http://localhost:8084/health")
	assert.Error(t, err) // 应该连接失败
}

func TestServerOptions(t *testing.T) {
	tests := []struct {
		name   string
		option ServerOptions
		check  func(*serverOptions) bool
	}{
		{
			name:   "WithName",
			option: WithName("custom-name"),
			check:  func(o *serverOptions) bool { return o.name == "custom-name" },
		},
		{
			name:   "WithAddress",
			option: WithAddress("0.0.0.0:9999"),
			check:  func(o *serverOptions) bool { return o.address == "0.0.0.0:9999" },
		},
		{
			name:   "WithMode",
			option: WithMode(gin.ReleaseMode),
			check:  func(o *serverOptions) bool { return o.mode == gin.ReleaseMode },
		},
		{
			name:   "WithTrans",
			option: WithTrans("en"),
			check:  func(o *serverOptions) bool { return o.trans == "en" },
		},
		{
			name:   "WithEnableProfiling",
			option: WithEnableProfiling(false),
			check:  func(o *serverOptions) bool { return o.enableProfiling == false },
		},
		{
			name:   "WithEnableMetrics",
			option: WithEnableMetrics(true),
			check:  func(o *serverOptions) bool { return o.enableMetrics == true },
		},
		{
			name:   "WithEnableTracing",
			option: WithEnableTracing(true),
			check:  func(o *serverOptions) bool { return o.enableTracing == true },
		},
		{
			name:   "WithEnableCert",
			option: WithEnableCert(true),
			check:  func(o *serverOptions) bool { return o.enableCert == true },
		},
		{
			name:   "WithCertFile",
			option: WithCertFile("/path/to/cert.pem"),
			check:  func(o *serverOptions) bool { return o.certFile == "/path/to/cert.pem" },
		},
		{
			name:   "WithKeyFile",
			option: WithKeyFile("/path/to/key.pem"),
			check:  func(o *serverOptions) bool { return o.keyFile == "/path/to/key.pem" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &serverOptions{
				name:            defaultName,
				address:         defaultAddr,
				mode:            gin.DebugMode,
				trans:           defaultrans,
				enableHealth:    true,
				enableProfiling: true,
				enableMetrics:   false,
				enableTracing:   false,
				enableCert:      false,
				certFile:        "",
				keyFile:         "",
				trustedProxies:  []string{},
				middlewares:     []gin.HandlerFunc{},
			}

			tt.option(opts)
			assert.True(t, tt.check(opts))
		})
	}
}

func TestServerCustomRoutes(t *testing.T) {
	server := NewServer(
		WithAddress("localhost:8085"),
	)

	// 添加自定义路由
	server.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test endpoint"})
	})

	server.POST("/echo", func(c *gin.Context) {
		var data map[string]interface{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, data)
	})

	// 启动服务器
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := server.Start(ctx)
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("Server start failed: %v", err)
		}
	}()

	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)

	// 测试 GET 路由
	resp, err := http.Get("http://localhost:8085/test")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 停止服务器
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopCancel()
	err = server.Stop(stopCtx)
	assert.NoError(t, err)
}

// 基准测试
func BenchmarkNewServer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewServer()
	}
}

func BenchmarkServerWithOptions(b *testing.B) {
	opts := []ServerOptions{
		WithName("bench-server"),
		WithAddress("localhost:8086"),
		WithMode(gin.ReleaseMode),
		WithEnableMetrics(true),
	}

	for i := 0; i < b.N; i++ {
		_ = NewServer(opts...)
	}
}
