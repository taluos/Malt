package rpcclient

import (
	"context"
	"crypto/tls"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/test/bufconn"
)

// 创建一个内存中的 gRPC 服务器用于测试
func newTestServer(t *testing.T) (*grpc.Server, *bufconn.Listener) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()

	t.Cleanup(func() {
		server.Stop()
		listener.Close()
	})

	go func() {
		if err := server.Serve(listener); err != nil {
			t.Errorf("服务器启动失败: %v", err)
		}
	}()

	return server, listener
}

// 创建一个自定义的 dialer 用于连接内存中的服务器
func getBufDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		return listener.Dial()
	}
}

// 生成测试用的TLS证书
func generateTestCert() (credentials.TransportCredentials, error) {

	return credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	}), nil
}

// 测试创建新客户端
func TestNewClient(t *testing.T) {
	// 创建测试服务器
	_, listener := newTestServer(t)

	// 创建自定义 dialer
	dialer := getBufDialer(listener)

	// 测试用例：使用默认选项创建客户端
	t.Run("使用默认选项", func(t *testing.T) {

		// 生成测试证书
		creds, err := generateTestCert()
		if err != nil {
			t.Fatalf("生成测试证书失败: %v", err)
		}

		client, err := NewClient(
			WithEndpoint("bufnet"),
			WithOptions(
				grpc.WithContextDialer(dialer),
				grpc.WithTransportCredentials(creds),
			),
		)

		if err != nil {
			t.Fatalf("创建客户端失败: %v", err)
		}

		if client == nil {
			t.Fatal("客户端不应为 nil")
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := client.Close(ctx); err != nil {
			t.Fatalf("关闭客户端失败: %v", err)
		}
	})

	// 测试用例：使用自定义选项创建客户端
	t.Run("使用自定义选项", func(t *testing.T) {
		// 生成测试证书
		creds, err := generateTestCert()
		if err != nil {
			t.Fatalf("生成测试证书失败: %v", err)
		}
		client, err := NewClient(
			WithEndpoint("bufnet"),
			WithTimeout(time.Second*10),
			WithInsecure(true),
			WithBalancerName("round_robin"),
			WithOptions(grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(creds)),
		)

		if err != nil {
			t.Fatalf("创建客户端失败: %v", err)
		}

		if client == nil {
			t.Fatal("客户端不应为 nil")
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := client.Close(ctx); err != nil {
			t.Fatalf("关闭客户端失败: %v", err)
		}
	})
}

// 测试获取客户端端点
func TestGetClientEndpoint(t *testing.T) {
	_, listener := newTestServer(t)
	dialer := getBufDialer(listener)
	// 生成测试证书
	creds, err := generateTestCert()
	if err != nil {
		t.Fatalf("生成测试证书失败: %v", err)
	}

	endpoint := "test-endpoint"
	client, err := NewClient(
		WithEndpoint(endpoint),
		WithOptions(grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(creds)),
	)

	if err != nil {
		t.Fatalf("创建客户端失败: %v", err)
	}

	if got := client.Endpoint(); got != endpoint {
		t.Errorf("GetClientEndpoint() = %v, 期望 %v", got, endpoint)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	client.Close(ctx)
}

// 测试客户端关闭
func TestClientClose(t *testing.T) {
	_, listener := newTestServer(t)
	dialer := getBufDialer(listener)

	// 测试正常关闭
	t.Run("正常关闭", func(t *testing.T) {
		// 生成测试证书
		creds, err := generateTestCert()
		if err != nil {
			t.Fatalf("生成测试证书失败: %v", err)
		}
		client, err := NewClient(
			WithEndpoint("bufnet"),
			WithOptions(grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(creds)),
		)

		if err != nil {
			t.Fatalf("创建客户端失败: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := client.Close(ctx); err != nil {
			t.Errorf("关闭客户端失败: %v", err)
		}
	})

	// 测试超时关闭
	t.Run("超时关闭", func(t *testing.T) {
		// 生成测试证书
		creds, err := generateTestCert()
		if err != nil {
			t.Fatalf("生成测试证书失败: %v", err)
		}
		client, err := NewClient(
			WithEndpoint("bufnet"),
			WithOptions(grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(creds)),
		)

		if err != nil {
			t.Fatalf("创建客户端失败: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
		defer cancel()

		// 这里可能会超时，取决于关闭速度
		_ = client.Close(ctx)
	})

	// 测试 nil 客户端连接
	t.Run("nil 连接", func(t *testing.T) {
		client := &Client{}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := client.Close(ctx); err != nil {
			t.Errorf("关闭 nil 连接应该不返回错误，但得到: %v", err)
		}
	})
}

// 测试 不安全连接
func TestDial(t *testing.T) {
	_, listener := newTestServer(t)
	dialer := getBufDialer(listener)

	// 测试安全连接
	t.Run("不安全连接", func(t *testing.T) {

		opts := clientOptions{
			endpoint: "bufnet",
			insecure: true,
			timeout:  time.Second,
			grpcOpts: []grpc.DialOption{grpc.WithContextDialer(dialer)},
		}

		conn, err := dial(true, opts)
		if err != nil {
			t.Fatalf("创建安全连接失败: %v", err)
		}

		if conn == nil {
			t.Fatal("连接不应为 nil")
		}

		conn.Close()
	})

	// 测试不安全连接
	t.Run("不安全连接", func(t *testing.T) {
		opts := clientOptions{
			endpoint: "bufnet",
			insecure: true,
			timeout:  time.Second,
			grpcOpts: []grpc.DialOption{grpc.WithContextDialer(dialer)},
		}

		conn, err := dial(true, opts)
		if err != nil {
			t.Fatalf("创建不安全连接失败: %v", err)
		}

		if conn == nil {
			t.Fatal("连接不应为 nil")
		}

		conn.Close()
	})
}
