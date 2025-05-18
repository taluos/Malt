package direct

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
)

// 测试 NewDirectResolver 函数
func TestNewDirectResolver(t *testing.T) {
	r := NewDirectResolver()
	assert.NotNil(t, r, "Resolver 不应该为 nil")
	assert.IsType(t, &directResolver{}, r, "Resolver 应该是 *directResolver 类型")
}

// 测试 directResolver 的 Close 方法
func TestDirectResolverClose(t *testing.T) {
	// 创建 mock ClientConn
	mockCC := new(mockClientConn)

	// 创建 directResolver 实例
	r := &directResolver{cc: mockCC}

	// 调用 Close 方法
	r.Close()

	// 由于 Close 方法是空实现，这里只是确保它不会崩溃
	assert.True(t, true, "Close 方法不应该引起 panic")
}

// 测试 directResolver 的 ResolveNow 方法
func TestDirectResolverResolveNow(t *testing.T) {
	// 创建 mock ClientConn
	mockCC := new(mockClientConn)

	// 创建 directResolver 实例
	r := &directResolver{cc: mockCC}

	// 调用 ResolveNow 方法
	opts := resolver.ResolveNowOptions{}
	r.ResolveNow(opts)

	// 由于 ResolveNow 方法是空实现，这里只是确保它不会崩溃
	assert.True(t, true, "ResolveNow 方法不应该引起 panic")
}
